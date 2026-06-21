package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	GitHubOwner   string
	GitHubRepo    string
	GitHubBranch  string
	RemoteBaseDir string

	GitHubToken string

	WorkDir     string
	RoadmapsDir string
	OutputDir   string
	StateFile   string

	RoadmapIDsFile string

	MongoURI    string
	MongoDBName string

	HTTPTimeout     time.Duration
	PipelineTimeout time.Duration
	WorkerCount     int
}

func Load() (Config, error) {
	cfg := Config{
		GitHubOwner:   getenv("GITHUB_OWNER", "nilbuild"),
		GitHubRepo:    getenv("GITHUB_REPO", "developer-roadmap"),
		GitHubBranch:  getenv("GITHUB_BRANCH", "master"),
		RemoteBaseDir: getenv("GITHUB_ROADMAPS_PATH", "src/data/roadmaps"),

		GitHubToken: os.Getenv("GITHUB_TOKEN"),

		WorkDir:        getenv("WORKDIR", "workdir"),
		RoadmapsDir:    getenv("LOCAL_ROADMAPS_DIR", "workdir/roadmaps"),
		OutputDir:      getenv("OUTPUT_DIR", "workdir/output"),
		StateFile:      getenv("STATE_FILE", "state/state.json"),
		RoadmapIDsFile: getenv("ROADMAP_IDS_FILE", "state/roadmap_ids.json"),

		MongoURI:    os.Getenv("MONGODB_URI"),
		MongoDBName: getenv("MONGODB_DB", "roadmapsdb"),

		HTTPTimeout:     30 * time.Second,
		PipelineTimeout: 15 * time.Minute,
		WorkerCount:     6,
	}

	if cfg.MongoURI == "" {
		return cfg, fmt.Errorf("MONGODB_URI is required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
