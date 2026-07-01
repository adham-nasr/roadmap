package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PipelineStore struct {
	client         *dynamodb.Client
	runTable       string
	completedTable string
}

func NewPipelineStore(client *dynamodb.Client, runTable, completedTable string) *PipelineStore {
	return &PipelineStore{
		client:         client,
		runTable:       runTable,
		completedTable: completedTable,
	}
}

func (s *PipelineStore) CreateRun(ctx context.Context, runID string, total int) error {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.runTable),
		Item: map[string]types.AttributeValue{
			"runId":     &types.AttributeValueMemberS{Value: runID},
			"total":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", total)},
			"remaining": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", total)},
		},
	})
	return err
}

func (s *PipelineStore) MarkRoadmapCompleted(ctx context.Context, runID, roadmapName string) (bool, error) {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.completedTable),
		Item: map[string]types.AttributeValue{
			"runId":       &types.AttributeValueMemberS{Value: runID},
			"roadmapName": &types.AttributeValueMemberS{Value: roadmapName},
		},
		ConditionExpression: aws.String("attribute_not_exists(runId) AND attribute_not_exists(roadmapName)"),
	})
	if err != nil {
		var cc *types.ConditionalCheckFailedException
		if errors.As(err, &cc) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *PipelineStore) DecrementRemaining(ctx context.Context, runID string) (int, int, error) {
	out, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.runTable),
		Key: map[string]types.AttributeValue{
			"runId": &types.AttributeValueMemberS{Value: runID},
		},
		UpdateExpression:    aws.String("SET remaining = remaining - :inc"),
		ConditionExpression: aws.String("remaining > :zero"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc":  &types.AttributeValueMemberN{Value: "1"},
			":zero": &types.AttributeValueMemberN{Value: "0"},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return 0, 0, err
	}
	remainingStr := out.Attributes["remaining"].(*types.AttributeValueMemberN).Value
	totalStr := out.Attributes["total"].(*types.AttributeValueMemberN).Value
	remaining, _ := strconv.Atoi(remainingStr)
	total, _ := strconv.Atoi(totalStr)
	return remaining, total, nil
}