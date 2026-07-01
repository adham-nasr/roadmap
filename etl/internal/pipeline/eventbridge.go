package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// Event represents a custom event to be sent to EventBridge.
type Event struct {
	// Source is the source of the event (e.g., "etl.pipeline").
	Source string
	// DetailType describes the event type (e.g., "RunComplete").
	DetailType string
	// Detail is the payload (can be any struct or map that marshals to JSON).
	Detail interface{}
}

// Send puts the event to the default or specified event bus.
func Send(ctx context.Context, event Event, eventBusName string) error {
	if eventBusName == "" {
		eventBusName = "default"
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := eventbridge.NewFromConfig(cfg)

	detailBytes, err := json.Marshal(event.Detail)
	if err != nil {
		return fmt.Errorf("failed to marshal event detail: %w", err)
	}

	entry := types.PutEventsRequestEntry{
		EventBusName: aws.String(eventBusName),
		Source:       aws.String(event.Source),
		DetailType:   aws.String(event.DetailType),
		Detail:       aws.String(string(detailBytes)),
		Time:         aws.Time(time.Now()),
	}
	_, err = client.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{entry},
	})
	if err != nil {
		return fmt.Errorf("failed to put event: %w", err)
	}
	return nil
}