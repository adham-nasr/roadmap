package load

import (
	"context"
)

// LoadFromOutput reads the JSON from the reader and upserts into MongoDB.
func LoadFromOutput(ctx context.Context, reader OutputReader, mongoURI, dbName string) error {
	// 1. Read outputs
	roadmaps, err := reader.ReadRoadmaps(ctx)
	if err != nil {
		return err
	}
	topics, err := reader.ReadTopics(ctx)
	if err != nil {
		return err
	}

	// 2. Connect to MongoDB
	client, err := Connect(ctx, mongoURI)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)

	// 3. Ensure indexes (optional – can be done once, but we'll do it each time for simplicity)
	if err := EnsureIndexes(ctx, db); err != nil {
		return err
	}

	// 4. Upsert data
	if err := UpsertRoadmaps(ctx, db, roadmaps); err != nil {
		return err
	}
	if err := UpsertTopics(ctx, db, topics); err != nil {
		return err
	}

	return nil
}