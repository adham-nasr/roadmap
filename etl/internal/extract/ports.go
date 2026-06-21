package extract

import (
	"context"

)
// SyncStateStore handles the persistent fingerprint state.
type SyncStateStore interface {
	LoadSyncState(ctx context.Context) (*State, error)
	SaveSyncState(ctx context.Context, state *State) error
}

// // IDStore handles the roadmap ID mappings.
// type IDStore interface {
// 	LoadIDStore(ctx context.Context) (*RoadmapIDStore, error)
// 	SaveIDStore(ctx context.Context, store *RoadmapIDStore) error
// }

// RawFileStore handles storing downloaded roadmap files.
type RawFileStore interface {
	// SaveRoadmapDirectory uploads the contents of a local temp directory
	// to persistent storage, keyed by roadmapName.
	SaveRoadmapDirectory(ctx context.Context, roadmapName string, localDirPath string) error
}