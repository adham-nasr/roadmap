package extract

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ETL/internal/util"
)

func (c *Client) DownloadRoadmap(ctx context.Context, localBaseDir, remoteBase string, rr RoadmapRemote) error {
	tmpDir := filepath.Join(localBaseDir, ".tmp_"+rr.Name)
	finalDir := filepath.Join(localBaseDir, rr.Name)

	_ = os.RemoveAll(tmpDir)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}

	prefix := strings.TrimSuffix(remoteBase, "/") + "/" + rr.Name + "/"

	for _, f := range rr.Files {
		if f.Type != "blob" {
			continue
		}

		rel := strings.TrimPrefix(f.Path, prefix)
		if rel == f.Path {
			continue
		}

		url := fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/%s/%s",
			c.owner, c.repo, c.branch, f.Path,
		)

		data, err := c.doBytes(ctx, url)
		if err != nil {
			return fmt.Errorf("download %s: %w", f.Path, err)
		}

		dst := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return err
		}
	}

	return util.AtomicReplaceDir(tmpDir, finalDir)
}
