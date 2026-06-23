package transform

import (
	"context"
	"encoding/json"   // <-- must be encoding/json, not just json
	"fmt"
	"path"
	"sort"
	"strings"
)

// ProcessAll reads all roadmaps from the reader, assigns IDs, and writes outputs.
func ProcessAll(ctx context.Context, reader RoadmapReader, idStore IDStore, outputWriter OutputWriter) error {
	// 1. Load ID store
	ids, err := idStore.LoadIDStore(ctx)
	if err != nil {
		return fmt.Errorf("load ID store: %w", err)
	}

	// 2. List roadmaps
	names, err := reader.ListRoadmaps(ctx)
	if err != nil {
		return fmt.Errorf("list roadmaps: %w", err)
	}
	if len(names) == 0 {
		// No roadmaps – nothing to do, but still write empty outputs? Possibly skip.
		return nil
	}

	// 3. Process each roadmap
	var roadmaps []*Roadmap
	for _, name := range names {
		rm, err := loadRoadmapFromReader(ctx, reader, name, ids)
		if err != nil {
			// Log and skip (or return error)
			// For now, skip
			continue
		}
		roadmaps = append(roadmaps, rm)
	}

	// 4. Save ID store (with any new IDs created)
	if err := idStore.SaveIDStore(ctx, ids); err != nil {
		return fmt.Errorf("save ID store: %w", err)
	}

	// 5. Collect outputs and sort
	sort.Slice(roadmaps, func(i, j int) bool { return roadmaps[i].ID < roadmaps[j].ID })
	roadmapOutputs, topicOutputs := CollectOutputs(roadmaps)

	// 6. Write outputs
	if err := outputWriter.WriteRoadmaps(ctx, roadmapOutputs); err != nil {
		return fmt.Errorf("write roadmaps: %w", err)
	}
	if err := outputWriter.WriteTopics(ctx, topicOutputs); err != nil {
		return fmt.Errorf("write topics: %w", err)
	}

	return nil
}

// loadRoadmapFromReader loads a single roadmap using the reader.
func loadRoadmapFromReader(ctx context.Context, reader RoadmapReader, name string, ids *RoadmapIDStore) (*Roadmap, error) {
	// Get the roadmap ID (creates if new)
	roadmapID, err := ids.GetOrCreate(name)
	if err != nil {
		return nil, err
	}

	// Read roadmap.json
	jsonData, err := reader.ReadFile(ctx, name, name+".json")
	if err != nil {
		return nil, fmt.Errorf("read roadmap.json: %w", err)
	}
	// Parse JSON structure (nodes & edges)
	var rf roadmapFileJSON
	if err := json.Unmarshal(jsonData, &rf); err != nil {
		return nil, fmt.Errorf("parse roadmap.json: %w", err)
	}

	// Read content files (list them)
	contentFiles, err := reader.ListFiles(ctx, name, "content")
	if err != nil {
		return nil, fmt.Errorf("list content files: %w", err)
	}
	// Build map from topic ID to content file name
	contentMap := make(map[string]string)
	for _, rel := range contentFiles {
		base := path.Base(rel)
		if !strings.HasSuffix(base, ".md") {
			continue
		}
		// filename is like "frontend@abc123.md" – we need the ID after '@'
		parts := strings.Split(strings.TrimSuffix(base, ".md"), "@")
		if len(parts) < 2 {
			continue
		}
		id := parts[len(parts)-1]
		contentMap[id] = rel
	}

	// Build topics map
	topics := make(map[string]*Topic)
	for _, info := range extractTopicInfos(rf.Nodes) {
		relPath, ok := contentMap[info.ID]
		if !ok {
			continue // topic has no content file – skip
		}
		mdData, err := reader.ReadFile(ctx, name, relPath)
		if err != nil {
			continue // skip this topic
		}
		// Parse the markdown content (similar to ParseFile, but from bytes)
		topic, err := parseMarkdownBytes(mdData)
		if err != nil {
			continue
		}
		topic.Info = info
		topic.RoadmapID = roadmapID
		topics[info.ID] = topic
	}

	// Link topics using edges
	filterAndLink(topics, rf.Edges)

	return &Roadmap{
		ID:     roadmapID,
		Name:   name,
		Topics: topics,
	}, nil
}

// parseMarkdownBytes is the same as ParseFile but reads from a byte slice.
// We'll reuse the existing parseLines function after splitting into lines.
func parseMarkdownBytes(data []byte) (*Topic, error) {
	lines := strings.Split(string(data), "\n")
	return parseLines(lines)
}