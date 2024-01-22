package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const Region = "eu-west-1"

// GetDDBClient creates and returns DynamoDB client using default AWS config
func CreateDDBClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = Region
		return nil
	})
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}
