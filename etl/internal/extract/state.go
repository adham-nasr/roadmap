package extract

import (
	"encoding/json"
	"os"
	"time"
)

type RoadmapState struct {
	Fingerprint string    `json:"fingerprint"`
	SyncedAt    time.Time `json:"syncedAt"`
}

type State struct {
	Roadmaps map[string]RoadmapState `json:"roadmaps"`
}

func LoadState(path string) (*State, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Roadmaps: map[string]RoadmapState{}}, nil
		}
		return nil, err
	}

	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if s.Roadmaps == nil {
		s.Roadmaps = map[string]RoadmapState{}
	}
	return &s, nil
}

func SaveState(path string, s *State) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}