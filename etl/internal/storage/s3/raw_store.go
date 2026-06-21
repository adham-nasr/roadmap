package s3

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type RawStore struct {
	client    *s3.Client
	bucket    string
	basePrefix string // e.g. "roadmaps"
}

func NewRawStore(client *s3.Client, bucket, basePrefix string) *RawStore {
	return &RawStore{client: client, bucket: bucket, basePrefix: basePrefix}
}

func (r *RawStore) SaveRoadmapDirectory(ctx context.Context, roadmapName string, localDirPath string) error {
	uploader := manager.NewUploader(r.client)
	
	err := filepath.WalkDir(localDirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(localDirPath, path)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s/%s/%s", r.basePrefix, roadmapName, filepath.ToSlash(rel))
		
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		
		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(r.bucket),
			Key:    aws.String(key),
			Body:   f,
		})
		return err
	})
	return err
}