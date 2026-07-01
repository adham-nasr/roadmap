package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ETL/internal/util"
)

// Task represents a single file to download.
type Task struct {
	URL      string // full HTTP URL
	DestPath string // absolute local file path
}

// Parallel downloads a list of tasks concurrently with the given concurrency.
func Parallel(ctx context.Context, tasks []Task, concurrency int) error {
	if concurrency <= 0 {
		concurrency = 5
	}
	return util.RunBounded(tasks, concurrency, func(task Task) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, task.URL, nil)
		if err != nil {
			return err
		}
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("HTTP %d for %s", resp.StatusCode, task.URL)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(task.DestPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(task.DestPath, body, 0644)
	})
}