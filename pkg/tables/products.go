package tables

import (
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/ddb"
)

// Products table
type Products struct {
	ddb.DynamoDBWritable

	Id       string `dynamodbav:"id"`
	Name     string `dynamodbav:"name"`
	Quantity int    `dynamodbav:"quantity"`

	// Version field is used for implementing Optimistic locking
	Version int `dynamodbav:"version"`
}
