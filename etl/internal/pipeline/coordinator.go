package pipeline

import (
	"context"

	"ETL/internal/storage/dynamodb"
)

type Coordinator interface {
	CreateRun(ctx context.Context, runID string, total int) error
	MarkRoadmapProcessed(ctx context.Context, runID, roadmapName string) (bool, error)
	DecrementRemaining(ctx context.Context, runID string) (remaining int, total int, err error)
}

type DynamoCoordinator struct {
	store *dynamodb.PipelineStore
}

func NewDynamoCoordinator(store *dynamodb.PipelineStore) *DynamoCoordinator {
	return &DynamoCoordinator{store: store}
}

func (c *DynamoCoordinator) CreateRun(ctx context.Context, runID string, total int) error {
	return c.store.CreateRun(ctx, runID, total)
}

func (c *DynamoCoordinator) MarkRoadmapProcessed(ctx context.Context, runID, roadmapName string) (bool, error) {
	return c.store.MarkRoadmapCompleted(ctx, runID, roadmapName)
}

func (c *DynamoCoordinator) DecrementRemaining(ctx context.Context, runID string) (remaining int, total int, err error) {
	return c.store.DecrementRemaining(ctx, runID)
}