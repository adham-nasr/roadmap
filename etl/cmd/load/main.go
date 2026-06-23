package main

import (
	"context"
	"log"
	"os"

	"ETL/internal/load"
	"ETL/internal/storage/s3"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// InputEvent is the payload from the Step Functions (or direct invocation).
type InputEvent struct {
	OutputBucket string `json:"outputBucket"`
	RoadmapsKey  string `json:"roadmapsKey"`
	TopicsKey    string `json:"topicsKey"`
}

type OutputEvent struct {
	Status string `json:"status"`
}

func handler(ctx context.Context, event InputEvent) (OutputEvent, error) {
	// Env vars
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DB")
	if mongoURI == "" || mongoDBName == "" {
		log.Fatal("Missing MongoDB environment variables")
	}

	// S3 client
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
	log.Printf("INFO: received event %+v", event)	// Adapter
	reader := s3.NewS3OutputReader(s3Client, event.OutputBucket, event.RoadmapsKey, event.TopicsKey)

	// Run load
	if err := load.LoadFromOutput(ctx, reader, mongoURI, mongoDBName); err != nil {
		return OutputEvent{}, err
	}

	return OutputEvent{Status: "success"}, nil
}

func main() {
	lambda.Start(handler)
}