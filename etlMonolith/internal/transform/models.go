package transform

type Position struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type ResourceType string

type Resource struct {
	Type  ResourceType `json:"type" bson:"type"`
	Title string       `json:"title" bson:"title"`
	Link  string       `json:"link" bson:"link"`
}

type EdgeRelation string

const (
	RelationTopic    EdgeRelation = "topic"
	RelationSubtopic EdgeRelation = "subtopic"
)

type ChildTopic struct {
	TargetID string       `json:"targetId" bson:"targetId"`
	Relation EdgeRelation `json:"relation" bson:"relation"`
}

type TopicInfo struct {
	ID       string
	Type     string
	Label    string
	Position Position
}

type Topic struct {
	Name        string
	Description string
	Resources   []Resource
	Info        TopicInfo
	ChildTopics []ChildTopic
	RoadmapID   string
}

type Roadmap struct {
	ID     string
	Name   string
	Topics map[string]*Topic
}
