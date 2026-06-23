package load

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"ETL/internal/transform"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	return client, nil
}

func UpsertRoadmaps(ctx context.Context, db *mongo.Database, roadmaps []transform.RoadmapOutput) error {
	col := db.Collection("roadmaps")
	models := make([]mongo.WriteModel, 0, len(roadmaps))
	now := time.Now().UTC()

	for _, r := range roadmaps {
		oid, err := bson.ObjectIDFromHex(r.ID)
		if err != nil {
			return fmt.Errorf("invalid roadmap object id for roadmap %s: %w", r.Name, err)
		}

		update := bson.M{
			"$set": bson.M{
				"name": r.Name,
				"source": bson.M{
					"syncedAt": now,
				},
			},
			"$setOnInsert": bson.M{
				"_id": oid,
			},
		}

		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": oid}).
			SetUpdate(update).
			SetUpsert(true))
	}

	if len(models) == 0 {
		return nil
	}

	_, err := col.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return err
}

func UpsertTopics(ctx context.Context, db *mongo.Database, topics []transform.TopicOutput) error {
	col := db.Collection("topics")
	models := make([]mongo.WriteModel, 0, len(topics))
	now := time.Now().UTC()

	for _, t := range topics {
		roadmapOID, err := bson.ObjectIDFromHex(t.RoadmapID)
		if err != nil {
			return fmt.Errorf("invalid roadmap object id in topic %s: %w", t.TopicID, err)
		}

		update := bson.M{
			"$set": bson.M{
				"topicId":     t.TopicID,
				"name":        t.Name,
				"label":       t.Label,
				"description": t.Description,
				"type":        t.Type,
				"roadmapId":   roadmapOID,
				"position":    t.Position,
				"resources":   t.Resources,
				"childTopics": t.ChildTopics,
				"source": bson.M{
					"syncedAt": now,
				},
			},
		}

		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"topicId": t.TopicID}).
			SetUpdate(update).
			SetUpsert(true))
	}

	if len(models) == 0 {
		return nil
	}

	_, err := col.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return err
}
