package s3

import (
	"context"
	"encoding/json"
	"ETL/internal/transform"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3OutputReader struct {
	client      *s3.Client
	bucket      string
	roadmapsKey string
	topicsKey   string
}

func NewS3OutputReader(client *s3.Client, bucket, roadmapsKey, topicsKey string) *S3OutputReader {
	return &S3OutputReader{
		client:      client,
		bucket:      bucket,
		roadmapsKey: roadmapsKey,
		topicsKey:   topicsKey,
	}
}

func (r *S3OutputReader) ReadRoadmaps(ctx context.Context) ([]transform.RoadmapOutput, error) {
	resp, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.roadmapsKey),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var roadmaps []transform.RoadmapOutput
	if err := json.NewDecoder(resp.Body).Decode(&roadmaps); err != nil {
		return nil, err
	}
	return roadmaps, nil
}

func (r *S3OutputReader) ReadTopics(ctx context.Context) ([]transform.TopicOutput, error) {
	resp, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.topicsKey),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var topics []transform.TopicOutput
	if err := json.NewDecoder(resp.Body).Decode(&topics); err != nil {
		return nil, err
	}
	return topics, nil
}