package load

import (
	"context"
	"ETL/internal/transform"
)

// OutputReader reads the transformed JSON outputs from persistent storage.
type OutputReader interface {
	ReadRoadmaps(ctx context.Context) ([]transform.RoadmapOutput, error)
	ReadTopics(ctx context.Context) ([]transform.TopicOutput, error)
}