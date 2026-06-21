package extract

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"ETL/internal/util"
)

type SyncResult struct {
	Eligible []RoadmapRemote
	Changed  []RoadmapRemote
	Skipped  []RoadmapRemote
}

func SyncRoadmaps(
	ctx context.Context,
	client *Client,
	statePath string,
	localRoadmapsDir string,
	remoteBase string,
	workers int,
) (*SyncResult, error) {
	if err := util.EnsureDir(filepath.Dir(statePath)); err != nil {
		return nil, err
	}
	if err := util.EnsureDir(localRoadmapsDir); err != nil {
		return nil, err
	}

	state, err := LoadState(statePath)
	if err != nil {
		return nil, err
	}

	tree, err := client.FetchTree(ctx)
	if err != nil {
		return nil, err
	}
	if tree.Truncated {
		return nil, fmt.Errorf("github tree response is truncated; fallback subtree strategy required")
	}

	eligible, err := DiscoverEligibleRoadmaps(tree, remoteBase)
	if err != nil {
		return nil, err
	}

	var changed []RoadmapRemote
	var skipped []RoadmapRemote

	for _, rr := range eligible {
		prev, ok := state.Roadmaps[rr.Name]
		if ok && prev.Fingerprint == rr.Fingerprint {
			skipped = append(skipped, rr)
			continue
		}
		changed = append(changed, rr)
	}

	err = util.RunBounded(changed, workers, func(rr RoadmapRemote) error {
		return client.DownloadRoadmap(ctx, localRoadmapsDir, remoteBase, rr)
	})
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	for _, rr := range changed {
		state.Roadmaps[rr.Name] = RoadmapState{
			Fingerprint: rr.Fingerprint,
			SyncedAt:    now,
		}
	}

	if err := SaveState(statePath, state); err != nil {
		return nil, err
	}

	return &SyncResult{
		Eligible: eligible,
		Changed:  changed,
		Skipped:  skipped,
	}, nil
}
