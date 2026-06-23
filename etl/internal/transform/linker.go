package transform

import "strings"

func filterAndLink(topics map[string]*Topic, edges []rawEdge) {
	for _, e := range edges {
		if e.Source == "" || e.Target == "" {
			continue
		}
		if _, ok := topics[e.Source]; !ok {
			continue
		}
		if _, ok := topics[e.Target]; !ok {
			continue
		}

		relation := RelationTopic
		if strings.EqualFold(e.Data.EdgeStyle, "dashed") {
			relation = RelationSubtopic
		}

		topics[e.Source].ChildTopics = append(topics[e.Source].ChildTopics, ChildTopic{
			TargetID: e.Target,
			Relation: relation,
		})
	}
}
