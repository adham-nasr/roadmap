package main

import (
	"context"
	"os"
    "log"
	"time"
	"ETL/internal/extract"
	"ETL/internal/storage/s3"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Response struct {
	RawBucket       string                   `json:"rawBucket"`
	ChangedRoadmaps []extract.RoadmapRemote  `json:"changedRoadmaps"`
}

func handler(ctx context.Context) (Response, error) {
	// Env
	// stateBucket := os.Getenv("STATE_BUCKET")
	rawBucket := os.Getenv("RAW_BUCKET")
	stateKey := os.Getenv("STATE_FILE_KEY") // e.g., "sync/state.json"
    _ = os.Getenv("IDS_FILE_KEY") // unused currently

	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	githubBranch := os.Getenv("GITHUB_BRANCH")
	githubToken := os.Getenv("GITHUB_TOKEN")
	remoteBase := os.Getenv("GITHUB_ROADMAPS_PATH")

	if rawBucket == "" || stateKey == "" {
		log.Printf("ERROR: Missing required S3 environment variables") 
		return Response{}, nil // Return error
	}
	if githubOwner == "" || githubRepo == "" || githubBranch == "" {
		log.Printf("ERROR: Missing required GitHub environment variables")
		return Response{}, nil // Return error
	}

	// Debuging starts
	if deadline, ok := ctx.Deadline(); ok {
		log.Printf("Incoming context deadline: %v, remaining: %v", deadline, time.Until(deadline))
	} else {
		log.Printf("Incoming context has no deadline")
	}
	// Debugging ends
	// AWS SDK
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
		return Response{}, err
	}

	endpoint := os.Getenv("AWS_ENDPOINT_URL_S3")

	var s3Client *awss3.Client

	if endpoint != "" {
		s3Client = awss3.NewFromConfig(
			cfg,
			func(o *awss3.Options) {
				o.BaseEndpoint = aws.String(endpoint)
				o.UsePathStyle = true
			},
		)
	} else {
		s3Client = awss3.NewFromConfig(cfg)
	}
	// Adapters
	log.Print("attempting to create adapters")
	syncStore := s3.NewSyncStore(s3Client, rawBucket, stateKey)
	rawStore := s3.NewRawStore(s3Client, rawBucket, "roadmaps")

	// GitHub Client
	log.Print("attempting to create github client")
	githubClient := extract.NewClient(githubOwner, githubRepo, githubBranch, githubToken, 30)
	log.Print("github client created")
	// Run core logic

	// Try then remve later start
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	// defer cancel()

	log.Print("attempting to run core logic")
	result, err := extract.SyncRoadmaps(ctx, githubClient, syncStore, rawStore, remoteBase, 18)
	if err != nil {
		return Response{}, err
	}

	_ = result    // Annoying Go

	return Response{
		RawBucket:       rawBucket,
		// ChangedRoadmaps: result.Changed,     FOR NOW ...............
	}, nil
}

func main() {
	log.Println("Starting  Lambda")
	lambda.Start(handler)
}