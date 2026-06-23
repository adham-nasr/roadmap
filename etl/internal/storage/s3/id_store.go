package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"ETL/internal/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3IDStore struct {
	client *s3.Client
	bucket string
	key    string // e.g. "sync/roadmap_ids.json"
}

func NewS3IDStore(client *s3.Client, bucket, key string) *S3IDStore {
	return &S3IDStore{client: client, bucket: bucket, key: key}
}

func (s *S3IDStore) LoadIDStore(ctx context.Context) (*transform.RoadmapIDStore, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	})
	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, &notFound) {
			return &transform.RoadmapIDStore{IDs: map[string]string{}}, nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	var store transform.RoadmapIDStore
	if err := json.NewDecoder(resp.Body).Decode(&store); err != nil {
		return nil, err
	}
	if store.IDs == nil {
		store.IDs = map[string]string{}
	}
	return &store, nil
}

func (s *S3IDStore) SaveIDStore(ctx context.Context, store *transform.RoadmapIDStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
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