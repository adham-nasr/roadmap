package transform

import (
	"encoding/json"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoadmapIDStore struct {
	IDs map[string]string `json:"ids"`
}

func LoadRoadmapIDStore(path string) (*RoadmapIDStore, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &RoadmapIDStore{IDs: map[string]string{}}, nil
		}
		return nil, err
	}

	var s RoadmapIDStore
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if s.IDs == nil {
		s.IDs = map[string]string{}
	}
	return &s, nil
}

func SaveRoadmapIDStore(path string, s *RoadmapIDStore) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func (s *RoadmapIDStore) GetOrCreate(name string) (string, error) {
	if id, ok := s.IDs[name]; ok {
		return id, nil
	}

	oid := bson.NewObjectID()
	hex := oid.Hex()
	s.IDs[name] = hex
	return hex, nil
}
