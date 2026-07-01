package extract

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"
	"os"
	"ETL/internal/util"
)

type SyncResult struct {
	Eligible []RoadmapRemote
	Changed  []RoadmapRemote
	Skipped  []RoadmapRemote
}

func SyncRoadmaps(
	ctx context.Context,
	github *Client,
	stateStore SyncStateStore,
	rawStore RawFileStore,
	remoteBase string,
	workers int,
) (*SyncResult, error) {
	// 1. Load state
	state, err := stateStore.LoadSyncState(ctx)
	if err != nil {
		return nil, err
	}
	if state.Roadmaps == nil {
		state.Roadmaps = map[string]RoadmapState{}
	}

	log.Print("fetching tree")
	// 2. Fetch tree from GitHub
	tree, err := github.FetchTree(ctx)
	if err != nil {
		log.Printf("error fetching tree: %v", err)
		return nil, err
	}
	if tree.Truncated {
		log.Print("github tree truncated")
		return nil, fmt.Errorf("github tree truncated")
	}


	log.Print("discovering eligible roadmaps")

	// 3. Discover eligible roadmaps
	eligible, err := DiscoverEligibleRoadmaps(tree, remoteBase)
	if err != nil {
		log.Printf("error discovering eligible roadmaps: %v", err)
		return nil, err
	}

	log.Printf("found %d eligible roadmaps",len(eligible))

	log.Printf("%+v",eligible)

	log.Print("--------------------------")

	return nil, nil
	log.Print("determining changed vs skipped")
	// 4. Determine changed vs skipped
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

	log.Print("downloading changed roadmaps")
	// 5. Download changed roadmaps to /tmp and upload to S3
	err = util.RunBounded(changed, workers, func(rr RoadmapRemote) error {
		// Download to a temporary directory
		tmpDir := filepath.Join("/tmp", ".tmp_"+rr.Name)
		if err := github.DownloadRoadmapToTemp(ctx, remoteBase, rr, tmpDir); err != nil {
			return err
		}
		// Upload the temp dir to S3
		if err := rawStore.SaveRoadmapDirectory(ctx, rr.Name, tmpDir); err != nil {
			return err
		}
		// Cleanup temp dir
		_ = os.RemoveAll(tmpDir)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 6. Update state and save
	now := time.Now().UTC()
	for _, rr := range changed {
		state.Roadmaps[rr.Name] = RoadmapState{
			Fingerprint: rr.Fingerprint,
			SyncedAt:    now,
		}
	}
	if err := stateStore.SaveSyncState(ctx, state); err != nil {
		return nil, err
	}

	return &SyncResult{
		Eligible: eligible,
		Changed:  changed,
		Skipped:  skipped,
	}, nil
}