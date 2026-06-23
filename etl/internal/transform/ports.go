package transform

import "context"

// RoadmapReader reads roadmap files from persistent storage (S3).
type RoadmapReader interface {
	// ListRoadmaps returns the names (directory names) of all available roadmaps.
	ListRoadmaps(ctx context.Context) ([]string, error)

	// ListFiles returns the relative file paths under a given subdirectory (e.g. "content").
	ListFiles(ctx context.Context, roadmapName, subDir string) ([]string, error)

	// ReadFile reads a specific file from a roadmap directory.
	ReadFile(ctx context.Context, roadmapName, relPath string) ([]byte, error)
}

// IDStore persists the mapping from roadmap name to its MongoDB ObjectID hex.
type IDStore interface {
	LoadIDStore(ctx context.Context) (*RoadmapIDStore, error)
	SaveIDStore(ctx context.Context, store *RoadmapIDStore) error
}

// OutputWriter writes the final JSON outputs.
type OutputWriter interface {
	WriteRoadmaps(ctx context.Context, roadmaps []RoadmapOutput) error
	WriteTopics(ctx context.Context, topics []TopicOutput) error
}