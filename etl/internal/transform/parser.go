package transform

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var resourceLineRegex = regexp.MustCompile(`^\s*-\s+\[@(\w+)@([^\]]+)\]\(([^)]+)\)`)

func ParseFile(filename string) (*Topic, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", filename, err)
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan %q: %w", filename, err)
	}

	return parseLines(lines)
}

func parseLines(lines []string) (*Topic, error) {
	type parserState int
	const (
		seekHeading parserState = iota
		collectDesc
		seekResources
		collectResources
	)

	topic := &Topic{
		Resources:   make([]Resource, 0),
		ChildTopics: make([]ChildTopic, 0),
	}

	state := seekHeading
	var descParts []string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		switch state {
		case seekHeading:
			if strings.HasPrefix(line, "# ") {
				topic.Name = strings.TrimPrefix(line, "# ")
				state = collectDesc
			}

		case collectDesc:
			if resourceLineRegex.MatchString(line) {
				state = collectResources
				appendResource(topic, line)
				continue
			}
			if line != "" {
				descParts = append(descParts, line)
			} else if len(descParts) > 0 {
				state = seekResources
			}

		case seekResources:
			if resourceLineRegex.MatchString(line) {
				state = collectResources
				appendResource(topic, line)
			}

		case collectResources:
			appendResource(topic, line)
		}
	}

	topic.Description = strings.Join(descParts, " ")
	return topic, nil
}

func appendResource(topic *Topic, line string) {
	m := resourceLineRegex.FindStringSubmatch(line)
	if m == nil {
		return
	}
	topic.Resources = append(topic.Resources, Resource{
		Type:  ResourceType(m[1]),
		Title: m[2],
		Link:  m[3],
	})
}
