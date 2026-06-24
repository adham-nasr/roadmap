package main

import (
	"context"
	"log"
	"os"

	"ETL/internal/storage/s3"
	"ETL/internal/transform"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type InputEvent struct {
	RawBucket string `json:"rawBucket"` // passed from Extract
	// we don't need changedRoadmaps; we'll process all.
}

type OutputEvent struct {
	OutputBucket string `json:"outputBucket"`
	RoadmapsKey  string `json:"roadmapsKey"`
	TopicsKey    string `json:"topicsKey"`
}

func handler(ctx context.Context, event InputEvent) (OutputEvent, error) {
	// Env
	// stateBucket := os.Getenv("STATE_BUCKET")
	outputBucket := os.Getenv("OUTPUT_BUCKET")
	stateIDsKey := os.Getenv("IDS_FILE_KEY") // e.g. "sync/roadmap_ids.json"
	roadmapsOutputKey := os.Getenv("ROADMAPS_OUTPUT_KEY") // e.g. "output/roadmaps.json"
	topicsOutputKey := os.Getenv("TOPICS_OUTPUT_KEY")     // e.g. "output/topics.json"

	if outputBucket == "" || stateIDsKey == "" {
		log.Fatal("Missing required S3 env vars")
	}

	// AWS SDK
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

	// Instantiate adapters
	reader := s3.NewS3RoadmapReader(s3Client, event.RawBucket, "roadmaps")
	idStore := s3.NewS3IDStore(s3Client, event.RawBucket, stateIDsKey)
	outputWriter := s3.NewS3OutputWriter(s3Client, outputBucket, roadmapsOutputKey, topicsOutputKey)

	// Run transform
	if err := transform.ProcessAll(ctx, reader, idStore, outputWriter); err != nil {
		return OutputEvent{}, err
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