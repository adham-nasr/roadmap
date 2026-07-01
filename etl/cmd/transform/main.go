package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"ETL/internal/pipeline"
	"ETL/internal/storage/s3"
	"ETL/internal/transform"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type OutputEvent struct {
	OutputBucket string `json:"outputBucket"`
	RoadmapsKey  string `json:"roadmapsKey"`
	TopicsKey    string `json:"topicsKey"`
}

func handler(ctx context.Context, event events.EventBridgeEvent) (OutputEvent, error) {
	// Parse detail to get rawBucket
	var detail struct {
		RawBucket string `json:"rawBucket"`
		RunID     string `json:"runId"`
	}
	if err := json.Unmarshal(event.Detail, &detail); err != nil {
		log.Printf("Failed to parse event detail: %v", err)
		return OutputEvent{}, err
	}
	rawBucket := detail.RawBucket

	// Env
	// stateBucket := os.Getenv("STATE_BUCKET")
	outputBucket := os.Getenv("OUTPUT_BUCKET")
	stateIDsKey := os.Getenv("IDS_FILE_KEY")
	roadmapsOutputKey := os.Getenv("ROADMAPS_OUTPUT_KEY")
	topicsOutputKey := os.Getenv("TOPICS_OUTPUT_KEY")
	eventBusName := os.Getenv("EVENT_BUS_NAME") // for emitting completion event

	if outputBucket == "" || stateIDsKey == "" || roadmapsOutputKey == "" || topicsOutputKey == "" {
		log.Fatal("Missing required S3 env vars")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return OutputEvent{}, err
	}
	endpoint := os.Getenv("AWS_ENDPOINT_URL_S3")
	var s3Client *awss3.Client
	if endpoint != "" {
		s3Client = awss3.NewFromConfig(cfg, func(o *awss3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		})
	} else {
		s3Client = awss3.NewFromConfig(cfg)
	}

	reader := s3.NewS3RoadmapReader(s3Client, rawBucket, "roadmaps")
	idStore := s3.NewS3IDStore(s3Client, rawBucket, stateIDsKey)
	outputWriter := s3.NewS3OutputWriter(s3Client, outputBucket, roadmapsOutputKey, topicsOutputKey)

	if err := transform.ProcessAll(ctx, reader, idStore, outputWriter); err != nil {
		return OutputEvent{}, err
	}

	// ✅ Emit EventBridge event to trigger Load Lambda
	log.Printf("Transform completed – emitting TransformComplete event")
	evt := pipeline.Event{
		Source:     "etl.pipeline",
		DetailType: "TransformComplete",
		Detail: map[string]string{
			"outputBucket": outputBucket,
			"roadmapsKey":  roadmapsOutputKey,
			"topicsKey":    topicsOutputKey,
		},
	}
	if err := pipeline.Send(ctx, evt, eventBusName); err != nil {
		log.Printf("Failed to emit TransformComplete event: %v", err)
		// Don't fail the Lambda – the transform succeeded, we can retry manually later.
	}

	return OutputEvent{
		OutputBucket: outputBucket,
		RoadmapsKey:  roadmapsOutputKey,
		TopicsKey:    topicsOutputKey,
	}, nil
}

func main() {
	lambda.Start(handler)
}