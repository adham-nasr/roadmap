package extract

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type TreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

type TreeResponse struct {
	SHA       string      `json:"sha"`
	Tree      []TreeEntry `json:"tree"`
	Truncated bool        `json:"truncated"`
}

type RoadmapRemote struct {
	Name        string
	Files       []TreeEntry
	Fingerprint string
}

func (c *Client) FetchTree(ctx context.Context) (*TreeResponse, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1",
		c.owner, c.repo, c.branch,
	)
	var tr TreeResponse
	if err := c.doJSON(ctx, url, &tr); err != nil {
		return nil, err
	}
	return &tr, nil
}

func DiscoverEligibleRoadmaps(tree *TreeResponse, remoteBase string) ([]RoadmapRemote, error) {
	prefix := strings.TrimSuffix(remoteBase, "/") + "/"
	grouped := map[string][]TreeEntry{}

	for _, e := range tree.Tree {
		if !strings.HasPrefix(e.Path, prefix) {
			continue
		}
		rest := strings.TrimPrefix(e.Path, prefix)
		parts := strings.Split(rest, "/")
		if len(parts) < 2 {
			continue
		}
		roadmap := parts[0]
		grouped[roadmap] = append(grouped[roadmap], e)
	}

	var remotes []RoadmapRemote
	for roadmap, files := range grouped {
		hasMigration := false
		hasMainJSON := false
		relevant := make([]TreeEntry, 0)

		for _, f := range files {
			rel := strings.TrimPrefix(f.Path, prefix+roadmap+"/")

			if rel == "migration-mapping.json" {
				hasMigration = true
			}
			if rel == roadmap+".json" {
				hasMainJSON = true
			}
			if rel == "migration-mapping.json" || rel == roadmap+".json" ||
				(strings.HasPrefix(rel, "content/") && strings.HasSuffix(rel, ".md")) {
				relevant = append(relevant, f)
			}
		}

		if !hasMigration || !hasMainJSON {
			continue
		}

		remotes = append(remotes, RoadmapRemote{
			Name:        roadmap,
			Files:       relevant,
			Fingerprint: fingerprint(relevant),
		})
	}

	sort.Slice(remotes, func(i, j int) bool { return remotes[i].Name < remotes[j].Name })
	return remotes, nil
}

func fingerprint(files []TreeEntry) string {
	rows := make([]string, 0, len(files))
	for _, f := range files {
		if f.Type != "blob" {
			continue
		}
		rows = append(rows, fmt.Sprintf("%s:%s", f.Path, f.SHA))
	}
	sort.Strings(rows)
	sum := sha256.Sum256([]byte(strings.Join(rows, "\n")))
	return hex.EncodeToString(sum[:])
}