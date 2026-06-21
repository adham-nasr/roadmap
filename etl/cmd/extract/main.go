package main

import (
	"context"
	"log"
	"os"

	"ETL/internal/extract"
	"ETL/internal/storage/s3"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Response struct {
	RawBucket       string                   `json:"rawBucket"`
	ChangedRoadmaps []extract.RoadmapRemote  `json:"changedRoadmaps"`
}

func handler(ctx context.Context) (Response, error) {
	// Env
	stateBucket := os.Getenv("STATE_BUCKET")
	rawBucket := os.Getenv("RAW_BUCKET")
	stateKey := os.Getenv("STATE_FILE_KEY") // e.g., "sync/state.json"
	idsKey := os.Getenv("IDS_FILE_KEY")     // e.g., "sync/roadmap_ids.json"
	
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	githubBranch := os.Getenv("GITHUB_BRANCH")
	githubToken := os.Getenv("GITHUB_TOKEN")
	remoteBase := os.Getenv("GITHUB_ROADMAPS_PATH")

	// AWS SDK
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return Response{}, err
	}
	s3Client := s3.NewFromConfig(cfg)

	// Adapters
	syncStore := s3.NewSyncStore(s3Client, stateBucket, stateKey)
	rawStore := s3.NewRawStore(s3Client, rawBucket, "roadmaps")

	// GitHub Client
	githubClient := extract.NewClient(githubOwner, githubRepo, githubBranch, githubToken, 30)

	// Run core logic
	result, err := extract.SyncRoadmaps(ctx, githubClient, syncStore, rawStore, remoteBase, 6)
	if err != nil {
		return Response{}, err
	}

	return Response{
		RawBucket:       rawBucket,
		ChangedRoadmaps: result.Changed,
	}, nil
}

func main() {
	lambda.Start(handler)
}