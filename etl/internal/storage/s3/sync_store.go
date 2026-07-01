package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"ETL/internal/extract"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// SyncStore is a factory that creates per‑roadmap state stores.
type SyncStore struct {
	client    *s3.Client
	bucket    string
	keyPrefix string // e.g., "sync/"
}

// NewSyncStore creates a new state store.
func NewSyncStore(client *s3.Client, bucket, keyPrefix string) *SyncStore {
	return &SyncStore{
		client:    client,
		bucket:    bucket,
		keyPrefix: keyPrefix,
	}
}

// ForRoadmap returns a scoped store for a specific roadmap.
func (s *SyncStore) ForRoadmap(name string) *RoadmapSyncStore {
	return &RoadmapSyncStore{
		client: s.client,
		bucket: s.bucket,
		key:    s.keyPrefix + name + ".json",
	}
}

// RoadmapSyncStore reads/writes state for a single roadmap.
type RoadmapSyncStore struct {
	client *s3.Client
	bucket string
	key    string
}

// Load reads the state for this roadmap. Returns nil if not found.
func (r *RoadmapSyncStore) Load(ctx context.Context) (*extract.RoadmapState, error) {
	resp, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.key),
	})
	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, &notFound) {
			return nil, nil // not found
		}
		return nil, err
	}
	defer resp.Body.Close()

	var state extract.RoadmapState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, err
	}
	return &state, nil
}

// Save stores the state for this roadmap.
func (r *RoadmapSyncStore) Save(ctx context.Context, state *extract.RoadmapState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.key),
		Body:   bytes.NewReader(data),
	})
	return err
}