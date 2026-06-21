package transform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type TopicOutput struct {
	TopicID     string       `json:"topicId" bson:"topicId"`
	Name        string       `json:"name" bson:"name"`
	Label       string       `json:"label" bson:"label"`
	Description string       `json:"description" bson:"description"`
	Type        string       `json:"type" bson:"type"`
	RoadmapID   string       `json:"roadmapId" bson:"roadmapId"`
	Position    Position     `json:"position" bson:"position"`
	Resources   []Resource   `json:"resources" bson:"resources"`
	ChildTopics []ChildTopic `json:"childTopics" bson:"childTopics"`
}

type RoadmapOutput struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func ToRoadmapOutput(r *Roadmap) RoadmapOutput {
	return RoadmapOutput{
		ID:   r.ID,
		Name: r.Name,
	}
}

func ToTopicOutput(t *Topic) TopicOutput {
	return TopicOutput{
		TopicID:     t.Info.ID,
		Name:        t.Name,
		Label:       t.Info.Label,
		Description: t.Description,
		Type:        t.Info.Type,
		RoadmapID:   t.RoadmapID,
		Position:    t.Info.Position,
		Resources:   t.Resources,
		ChildTopics: t.ChildTopics,
	}
}

func CollectOutputs(roadmaps []*Roadmap) ([]RoadmapOutput, []TopicOutput) {
	rms := make([]RoadmapOutput, 0, len(roadmaps))
	topics := make([]TopicOutput, 0)

	for _, r := range roadmaps {
		rms = append(rms, ToRoadmapOutput(r))
		for _, t := range r.Topics {
			topics = append(topics, ToTopicOutput(t))
		}
	}

	sort.Slice(rms, func(i, j int) bool { return rms[i].ID < rms[j].ID })
	sort.Slice(topics, func(i, j int) bool { return topics[i].TopicID < topics[j].TopicID })

	return rms, topics
}

func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
