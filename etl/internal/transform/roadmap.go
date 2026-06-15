package transform

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type node struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Position Position `json:"position"`
	Data     struct {
		Label string `json:"label"`
	} `json:"data"`
}

type rawEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Data   struct {
		EdgeStyle string `json:"edgeStyle"`
	} `json:"data"`
}

type roadmapFileJSON struct {
	Nodes []node    `json:"nodes"`
	Edges []rawEdge `json:"edges"`
}

func loadRoadmapFile(roadmapPath, roadmapName string) (*roadmapFileJSON, error) {
	filePath := filepath.Join(roadmapPath, roadmapName+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", filePath, err)
	}

	var rf roadmapFileJSON
	if err := json.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("parse %q: %w", filePath, err)
	}
	return &rf, nil
}

func extractTopicInfos(nodes []node) []TopicInfo {
	out := make([]TopicInfo, 0)
	for _, n := range nodes {
		if n.Type != "topic" && n.Type != "subtopic" {
			continue
		}
		out = append(out, TopicInfo{
			ID:       n.ID,
			Type:     n.Type,
			Label:    n.Data.Label,
			Position: n.Position,
		})
	}
	return out
}

func getContentIDs(contentDir string) (map[string]string, error) {
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}

	ids := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) != ".md" {
			continue
		}
		base := strings.TrimSuffix(name, ".md")
		parts := strings.Split(base, "@")
		if len(parts) < 2 {
			continue
		}
		id := parts[len(parts)-1]
		ids[id] = name
	}
	return ids, nil
}

func LoadRoadmap(name, folderPath string, idStore *RoadmapIDStore) (*Roadmap, error) {
	rf, err := loadRoadmapFile(folderPath, name)
	if err != nil {
		return nil, err
	}

	roadmapID, err := idStore.GetOrCreate(name)
	if err != nil {
		return nil, err
	}

	rm := &Roadmap{
		ID:     roadmapID,
		Name:   name,
		Topics: map[string]*Topic{},
	}

	contentIDs, err := getContentIDs(filepath.Join(folderPath, "content"))
	if err != nil {
		return nil, err
	}

	for _, info := range extractTopicInfos(rf.Nodes) {
		filename, ok := contentIDs[info.ID]
		if !ok {
			continue
		}

		topic, err := ParseFile(filepath.Join(folderPath, "content", filename))
		if err != nil {
			continue
		}

		topic.Info = info
		topic.RoadmapID = rm.ID
		rm.Topics[info.ID] = topic
	}

	filterAndLink(rm.Topics, rf.Edges)
	return rm, nil
}
