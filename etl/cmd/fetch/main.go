package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"ETL/internal/extract"
	"ETL/internal/pipeline"
	dynamodbstore "ETL/internal/storage/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSMessage struct {
	RunID       string   `json:"runId"`
	RoadmapName string   `json:"roadmapName"`
	Fingerprint string   `json:"fingerprint"`
	RawBucket   string   `json:"rawBucket"`
	BaseURL     string   `json:"baseURL"`
	Files       []string `json:"files"`
}

func handler(ctx context.Context) error {
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	githubBranch := os.Getenv("GITHUB_BRANCH")
	githubToken := os.Getenv("GITHUB_TOKEN")
	remoteBase := os.Getenv("GITHUB_ROADMAPS_PATH")
	queueURL := os.Getenv("SQS_QUEUE_URL")
	rawBucket := os.Getenv("RAW_BUCKET")
	runTable := os.Getenv("RUN_TABLE")
	completedTable := os.Getenv("COMPLETED_TABLE")

	if githubOwner == "" || githubRepo == "" || githubBranch == "" || remoteBase == "" || queueURL == "" || rawBucket == "" || runTable == "" || completedTable == "" {
		log.Fatal("Missing required environment variables")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	sqsClient := sqs.NewFromConfig(cfg)
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// Create pipeline store and coordinator
	pipelineStore := dynamodbstore.NewPipelineStore(dynamoClient, runTable, completedTable)
	coord := pipeline.NewDynamoCoordinator(pipelineStore)

	githubClient := extract.NewClient(githubOwner, githubRepo, githubBranch, githubToken, 30*time.Second)

	tree, err := githubClient.FetchTree(ctx)
	if err != nil {
		return err
	}
	if tree.Truncated {
		log.Println("Tree truncated")
		return nil
	}

	eligible, err := extract.DiscoverEligibleRoadmaps(tree, remoteBase)
	if err != nil {
		return err
	}
	if len(eligible) == 0 {
		log.Println("No eligible roadmaps found")
		return nil
	}

	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	if err := coord.CreateRun(ctx, runID, len(eligible)); err != nil {
		return err
	}
	log.Printf("Created run %s with %d roadmaps", runID, len(eligible))

	prefixBase := strings.TrimSuffix(remoteBase, "/") + "/"
	sentCount := 0

	for _, rr := range eligible {
		baseURL := fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/%s/%s/%s/",
			githubOwner, githubRepo, githubBranch, remoteBase, rr.Name,
		)
		var relPaths []string
		roadmapPrefix := prefixBase + rr.Name + "/"
		for _, entry := range rr.Files {
			if entry.Type != "blob" {
				continue
			}
			rel := strings.TrimPrefix(entry.Path, roadmapPrefix)
			if rel == entry.Path {
				continue
			}
			relPaths = append(relPaths, rel)
		}
		if len(relPaths) == 0 {
			continue
		}

		msg := SQSMessage{
			RunID:       runID,
			RoadmapName: rr.Name,
			Fingerprint: rr.Fingerprint,
			RawBucket:   rawBucket,
			BaseURL:     baseURL,
			Files:       relPaths,
		}
		body, _ := json.Marshal(msg)

		_, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String(string(body)),
		})
		if err != nil {
			log.Printf("Failed to send message for %s: %v", rr.Name, err)
			continue
		}
		sentCount++
	}
	log.Printf("Sent %d messages for run %s", sentCount, runID)
	return nil
}

func main() {
	lambda.Start(handler)
}