package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"ETL/internal/extract"
	// "ETL/internal/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type SyncStore struct {
	client *s3.Client
	bucket string
	key    string
}

func NewSyncStore(client *s3.Client, bucket, key string) *SyncStore {
	return &SyncStore{client: client, bucket: bucket, key: key}
}

func (s *SyncStore) LoadSyncState(ctx context.Context) (*extract.State, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	})
	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, &notFound) {
			return &extract.State{Roadmaps: map[string]extract.RoadmapState{}}, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	var state extract.State
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, err
	}
	if state.Roadmaps == nil {
		state.Roadmaps = map[string]extract.RoadmapState{}
	}
	return &state, nil
}

func (s *SyncStore) SaveSyncState(ctx context.Context, state *extract.State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
		Body:   bytes.NewReader(data),
	})
	return err
}