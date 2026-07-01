package load

import (
	"context"
	"log"
)

func LoadFromOutput(ctx context.Context, reader OutputReader, mongoURI, dbName string) error {
	log.Println("LoadFromOutput: reading roadmaps from S3...")
	roadmaps, err := reader.ReadRoadmaps(ctx)
	if err != nil {
		log.Printf("LoadFromOutput: failed to read roadmaps: %v", err)
		return err
	}
	log.Printf("LoadFromOutput: read %d roadmaps", len(roadmaps))

	log.Println("LoadFromOutput: reading topics from S3...")
	topics, err := reader.ReadTopics(ctx)
	if err != nil {
		log.Printf("LoadFromOutput: failed to read topics: %v", err)
		return err
	}
	log.Printf("LoadFromOutput: read %d topics", len(topics))

	log.Println("LoadFromOutput: connecting to MongoDB...")
	client, err := Connect(ctx, mongoURI)
	if err != nil {
		log.Printf("LoadFromOutput: failed to connect to MongoDB: %v", err)
		return err
	}
	defer client.Disconnect(ctx)
	log.Println("LoadFromOutput: connected to MongoDB")

	db := client.Database(dbName)

	log.Println("LoadFromOutput: ensuring indexes...")
	if err := EnsureIndexes(ctx, db); err != nil {
		log.Printf("LoadFromOutput: failed to ensure indexes: %v", err)
		return err
	}
	log.Println("LoadFromOutput: indexes ensured")

	log.Println("LoadFromOutput: upserting roadmaps...")
	if err := UpsertRoadmaps(ctx, db, roadmaps); err != nil {
		log.Printf("LoadFromOutput: failed to upsert roadmaps: %v", err)
		return err
	}
	log.Println("LoadFromOutput: roadmaps upserted")

	log.Println("LoadFromOutput: upserting topics...")
	if err := UpsertTopics(ctx, db, topics); err != nil {
		log.Printf("LoadFromOutput: failed to upsert topics: %v", err)
		return err
	}
	log.Println("LoadFromOutput: topics upserted")

	log.Println("LoadFromOutput: completed successfully")
	return nil
}