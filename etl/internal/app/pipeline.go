package app

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"time"

	"ETL/internal/config"
	"ETL/internal/extract"
	"ETL/internal/load"
	"ETL/internal/transform"
)

func Run(ctx context.Context, cfg config.Config) error {
	extractor := extract.NewClient(
		cfg.GitHubOwner,
		cfg.GitHubRepo,
		cfg.GitHubBranch,
		cfg.GitHubToken,
		cfg.HTTPTimeout,
	)

	syncResult, err := extract.SyncRoadmaps(
		ctx,
		extractor,
		cfg.StateFile,
		cfg.RoadmapsDir,
		cfg.RemoteBaseDir,
		cfg.WorkerCount,
	)
	if err != nil {
		return err
	}

	_ = syncResult

	idStore, err := transform.LoadRoadmapIDStore(cfg.RoadmapIDsFile)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(cfg.RoadmapsDir)
	if err != nil {
		return err
	}

	var roadmaps []*transform.Roadmap
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		rm, err := transform.LoadRoadmap(name, filepath.Join(cfg.RoadmapsDir, name), idStore)
		if err != nil {
			continue
		}
		roadmaps = append(roadmaps, rm)
	}

	sort.Slice(roadmaps, func(i, j int) bool { return roadmaps[i].ID < roadmaps[j].ID })

	if err := transform.SaveRoadmapIDStore(cfg.RoadmapIDsFile, idStore); err != nil {
		return err
	}

	roadmapOutputs, topicOutputs := transform.CollectOutputs(roadmaps)

	if err := transform.WriteJSON(filepath.Join(cfg.OutputDir, "roadmaps.json"), roadmapOutputs); err != nil {
		return err
	}
	if err := transform.WriteJSON(filepath.Join(cfg.OutputDir, "topics.json"), topicOutputs); err != nil {
		return err
	}

	client, err := load.Connect(ctx, cfg.MongoURI)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	db := client.Database(cfg.MongoDBName)

	loadCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	if err := load.EnsureIndexes(loadCtx, db); err != nil {
		return err
	}
	if err := load.UpsertRoadmaps(loadCtx, db, roadmapOutputs); err != nil {
		return err
	}
	if err := load.UpsertTopics(loadCtx, db, topicOutputs); err != nil {
		return err
	}

	return nil
}
