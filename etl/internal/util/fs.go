package util

import (
	"os"
	"path/filepath"
)

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func AtomicReplaceDir(tmpDir, finalDir string) error {
	_ = os.RemoveAll(finalDir)
	parent := filepath.Dir(finalDir)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}
	return os.Rename(tmpDir, finalDir)
}