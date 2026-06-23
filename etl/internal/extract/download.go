package extract

import (
	"context"
	"fmt"
	"os"
	"log"
	"path/filepath"
	"strings"
)

func (c *Client) DownloadRoadmapToTemp(ctx context.Context, remoteBase string, rr RoadmapRemote, tmpDir string) error {
	_ = os.RemoveAll(tmpDir)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}

	prefix := strings.TrimSuffix(remoteBase, "/") + "/" + rr.Name + "/"
	log.Printf("Downloading roadmap %s with %d files", rr.Name, len(rr.Files))
	for i, f := range rr.Files {
		if f.Type != "blob" {
			continue
		}
		rel := strings.TrimPrefix(f.Path, prefix)
		if rel == f.Path {
			continue
		}
		if i%10 == 0 {
            log.Printf("  %s: downloaded %d/%d files", rr.Name, i, len(rr.Files))
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
	log.Printf("Finished downloading roadmap %s", rr.Name)
	return nil
}