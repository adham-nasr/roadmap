package load

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	topics := db.Collection("topics")

	_, err := topics.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "topicId", Value: 1},
			},
			Options: options.Index().
				SetName("ux_topics_topicId").
				SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "roadmapId", Value: 1},
			},
			Options: options.Index().
				SetName("idx_topics_roadmapId"),
		},
		{
			Keys: bson.D{
				{Key: "roadmapId", Value: 1},
				{Key: "type", Value: 1},
			},
			Options: options.Index().
				SetName("idx_topics_roadmapId_type"),
		},
	})
	return err
}
