package ddb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBWritable is a type which can be serialized to a DynamoDB table
type DynamoDBWritable interface{}

// DynamoDB is a wrapper around AWS SDK's client to unify the way of using
// SDK's client and also abstract complexities involved in using the main library.
type DynamoDB struct {
	DynamoDbClient *dynamodb.Client
}

func NewDynamoDBClient() *DynamoDB {
	return &DynamoDB{DynamoDbClient: CreateDDBClient()}
}

// AddBatch writes a batch of new records to table respecting the given max number
func (d DynamoDB) AddBatch(tableName string, records []DynamoDBWritable, max int) (int, error) {
	var err error
	var item map[string]types.AttributeValue
	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	for start < max && start < len(records) {
		var writeReqs []types.WriteRequest
		if end > len(records) {
			end = len(records)
		}
		for _, movie := range records[start:end] {
			item, err = attributevalue.MarshalMap(movie)

			if err != nil {
				log.Printf("Couldn't marshal record for batch writing: %v\n", err)
			} else {
				writeReqs = append(
					writeReqs,
					types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
				)
			}
		}
		_, err = d.DynamoDbClient.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{tableName: writeReqs}})
		if err != nil {
			log.Printf("Couldn't add a batch of records to %v: %v\n", tableName, err)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	return written, err
}

// GetItem returns the requested key from the given table
func (d DynamoDB) GetItem(tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	result, err := d.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       key,
	})
	if err != nil {
		log.Printf("Couldn't get item from %v: %v\n", tableName, err)
		return nil, err
	}
	return result.Item, nil
}
