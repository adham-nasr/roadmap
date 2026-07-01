package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"ETL/internal/download"
	"ETL/internal/extract"
	"ETL/internal/pipeline"
	dynamodbstore "ETL/internal/storage/dynamodb"
	"ETL/internal/storage/s3"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// SQSMessage matches the dispatcher's payload.
type SQSMessage struct {
	RunID       string   `json:"runId"`
	RoadmapName string   `json:"roadmapName"`
	Fingerprint string   `json:"fingerprint"`
	RawBucket   string   `json:"rawBucket"`
	BaseURL     string   `json:"baseURL"`
	Files       []string `json:"files"`
}

// handler is the Lambda entry point.
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Environment
	stateBucket := os.Getenv("RAW_BUCKET")
	stateKeyPrefix := os.Getenv("STATE_KEY_PREFIX")
	runTable := os.Getenv("RUN_TABLE")
	completedTable := os.Getenv("COMPLETED_TABLE")
	eventBusName := os.Getenv("EVENT_BUS_NAME")

	if stateBucket == "" || stateKeyPrefix == "" || runTable == "" || completedTable == "" {
		log.Fatal("Missing required env vars: STATE_BUCKET, STATE_KEY_PREFIX, RUN_TABLE, COMPLETED_TABLE")
	}

	// AWS SDK
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	s3Client := awss3.NewFromConfig(cfg)
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// State store (S3 per‑roadmap)
	stateStore := s3.NewSyncStore(s3Client, stateBucket, stateKeyPrefix)

	// Pipeline coordinator (DynamoDB)
	pipelineStore := dynamodbstore.NewPipelineStore(dynamoClient, runTable, completedTable)
	coord := pipeline.NewDynamoCoordinator(pipelineStore)

	// Process each SQS message
	for _, record := range sqsEvent.Records {
		var msg SQSMessage
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue // skip and move to DLQ after maxReceiveCount
		}

		log.Printf("Processing roadmap: %s (run: %s)", msg.RoadmapName, msg.RunID)

		// 1. Deduplicate (insert into CompletedRoadmaps if first time)
		firstTime, err := coord.MarkRoadmapProcessed(ctx, msg.RunID, msg.RoadmapName)
		if err != nil {
			return err
		}
		if !firstTime {
			log.Printf("Duplicate message for %s – skipping", msg.RoadmapName)
			continue
		}

		// 2. Check if already up‑to‑date (state check)
		roadmapStore := stateStore.ForRoadmap(msg.RoadmapName)
		curr, err := roadmapStore.Load(ctx)
		if err != nil {
			return err
		}
		if curr != nil && curr.Fingerprint == msg.Fingerprint {
			log.Printf("Skipping %s – already up‑to‑date", msg.RoadmapName)
			continue
		}

		// 3. Download files to /tmp
		tmpDir := filepath.Join("/tmp", ".tmp_"+msg.RoadmapName)
		if err := os.RemoveAll(tmpDir); err != nil {
			return err
		}
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return err
		}

		tasks := make([]download.Task, len(msg.Files))
		for i, relPath := range msg.Files {
			tasks[i] = download.Task{
				URL:      msg.BaseURL + relPath,
				DestPath: filepath.Join(tmpDir, filepath.FromSlash(relPath)),
			}
		}
		if err := download.Parallel(ctx, tasks, 10); err != nil {
			_ = os.RemoveAll(tmpDir)
			return err
		}

		// 4. Upload to raw bucket
		rawStore := s3.NewRawStore(s3Client, msg.RawBucket, "roadmaps")
		if err := rawStore.SaveRoadmapDirectory(ctx, msg.RoadmapName, tmpDir); err != nil {
			_ = os.RemoveAll(tmpDir)
			return err
		}

		// 5. Update state (S3 per‑roadmap)
		if err := roadmapStore.Save(ctx, &extract.RoadmapState{
			Fingerprint: msg.Fingerprint,
			SyncedAt:    time.Now().UTC(),
		}); err != nil {
			_ = os.RemoveAll(tmpDir)
			return err
		}

		// 6. Clean up temporary directory
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("Warning: cleanup failed: %v", err)
		}

		// 7. Atomic decrement of remaining counter
		remaining, total, err := coord.DecrementRemaining(ctx, msg.RunID)
		if err != nil {
			log.Printf("Failed to decrement remaining for %s: %v", msg.RoadmapName, err)
			return err
		}
		log.Printf("Remaining: %d/%d for run %s", remaining, total, msg.RunID)

		// 8. If this was the last roadmap, trigger Transform Lambda
		if remaining == 0 {
			log.Printf("All roadmaps completed for run %s – emitting completion event", msg.RunID)
			event := pipeline.Event{
				Source:     "etl.pipeline",
				DetailType: "RunComplete",
				Detail: map[string]string{
					"rawBucket": msg.RawBucket,
					"runId":     msg.RunID,
				},
			}
			if err := pipeline.Send(ctx, event, eventBusName); err != nil {
				log.Printf("Failed to emit completion event: %v", err)
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}