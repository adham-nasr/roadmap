package s3

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3RoadmapReader struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewS3RoadmapReader(client *s3.Client, bucket, prefix string) *S3RoadmapReader {
	return &S3RoadmapReader{client: client, bucket: bucket, prefix: prefix}
}

func (r *S3RoadmapReader) ListRoadmaps(ctx context.Context) ([]string, error) {
	_ = types.CommonPrefix{} // dummy usage to silence "unused import"

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(r.bucket),
		Prefix:    aws.String(r.prefix + "/"),
		Delimiter: aws.String("/"),
	}
	resp, err := r.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, common := range resp.CommonPrefixes {
		prefixStr := aws.ToString(common.Prefix)
		trimmed := strings.TrimPrefix(prefixStr, r.prefix+"/")
		trimmed = strings.TrimSuffix(trimmed, "/")
		if trimmed != "" {
			names = append(names, trimmed)
		}
	}
	return names, nil
}

func (r *S3RoadmapReader) ListFiles(ctx context.Context, roadmapName, subDir string) ([]string, error) {
	basePrefix := path.Join(r.prefix, roadmapName, subDir) + "/"
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(basePrefix),
	}
	resp, err := r.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}
	var relPaths []string
	for _, obj := range resp.Contents {
		key := aws.ToString(obj.Key)
		rel := strings.TrimPrefix(key, path.Join(r.prefix, roadmapName)+"/")
		relPaths = append(relPaths, rel)
	}
	return relPaths, nil
}

func (r *S3RoadmapReader) ReadFile(ctx context.Context, roadmapName, relPath string) ([]byte, error) {
	key := path.Join(r.prefix, roadmapName, relPath)
	resp, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}