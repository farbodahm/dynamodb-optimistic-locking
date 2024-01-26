package tables

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/ddb"
)

// Products table
type Product struct {
	ddb.DynamoDBWritable `dynamodbav:",omitempty"`

	Id       string `dynamodbav:"id"`
	Name     string `dynamodbav:"name"`
	Quantity int    `dynamodbav:"quantity"`

	// Version field is used for implementing Optimistic locking
	Version int `dynamodbav:"version"`
}

// GetProduct returns a product with the given id from the products table.
// If no product with the given id exists, an empty product is returned.
func GetProduct(d ddb.DynamoDB, id string) (Product, error) {
	var product Product

	marshaledId, err := attributevalue.Marshal(id)

	if err != nil {
		log.Printf("WARN: Couldn't marshal product id %v\n", id)
		return product, err
	}

	query := map[string]types.AttributeValue{"id": marshaledId}
	result, err := d.GetItem("products", query)

	if err != nil {
		log.Printf("WARN: Couldn't get product id %v\n", id)
		return product, err
	}

	if err := attributevalue.UnmarshalMap(result, &product); err != nil {
		log.Printf("WARN: Couldn't unmarshal product id %v\n", id)
		return product, err
	}

	return product, nil
}
