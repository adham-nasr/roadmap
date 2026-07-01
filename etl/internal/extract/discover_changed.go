package extract

import (
	"context"
	"fmt"
)

func DiscoverChangedRoadmaps(ctx context.Context, client *Client, stateStore SyncStateStore, remoteBase string) ([]RoadmapRemote, error) {
	state, err := stateStore.LoadSyncState(ctx)
	if err != nil {
		return nil, err
	}
	if state.Roadmaps == nil {
		state.Roadmaps = map[string]RoadmapState{}
	}

	tree, err := client.FetchTree(ctx)
	if err != nil {
		return nil, err
	}
	if tree.Truncated {
		return nil, fmt.Errorf("tree truncated")
	}

	eligible, err := DiscoverEligibleRoadmaps(tree, remoteBase)
	if err != nil {
		return nil, err
	}

	var changed []RoadmapRemote
	for _, rr := range eligible {
		prev, ok := state.Roadmaps[rr.Name]
		if !ok || prev.Fingerprint != rr.Fingerprint {
			changed = append(changed, rr)
		}
	}
	return changed, nil
}