package s3

import (
	"bytes"
	"context"
	"encoding/json"

	"ETL/internal/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3OutputWriter struct {
	client      *s3.Client
	bucket      string
	roadmapsKey string // e.g. "output/roadmaps.json"
	topicsKey   string // e.g. "output/topics.json"
}

func NewS3OutputWriter(client *s3.Client, bucket, roadmapsKey, topicsKey string) *S3OutputWriter {
	return &S3OutputWriter{
		client:      client,
		bucket:      bucket,
		roadmapsKey: roadmapsKey,
		topicsKey:   topicsKey,
	}
}

func (w *S3OutputWriter) WriteRoadmaps(ctx context.Context, roadmaps []transform.RoadmapOutput) error {
	data, err := json.MarshalIndent(roadmaps, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(w.roadmapsKey),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (w *S3OutputWriter) WriteTopics(ctx context.Context, topics []transform.TopicOutput) error {
	data, err := json.MarshalIndent(topics, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(w.topicsKey),
		Body:   bytes.NewReader(data),
	})
	return err
}