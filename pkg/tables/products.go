package tables

import (
	"fmt"
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
func GetProduct(d ddb.DynamoDB, tableName string, id string) (Product, error) {
	var product Product

	marshaledId, err := attributevalue.Marshal(id)

	if err != nil {
		log.Printf("WARN: Couldn't marshal product id %v\n", id)
		return product, err
	}

	query := map[string]types.AttributeValue{"id": marshaledId}
	result, err := d.GetItem(tableName, query)

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

// SafeUpdateProductQuantity updates the quantity of a Product in the products table with thread safety.
// It employs an optimistic locking mechanism to efficiently manage concurrent updates.
func SafeUpdateProductQuantity(d ddb.DynamoDB, tableName string, product Product) error {
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: product.Id},
	}
	updateExpression := "SET quantity = :newQuantity, version = :newVersion"
	conditionExpression := "version = :oldVersion"

	newQuantity := fmt.Sprint(product.Quantity)
	oldVersion := fmt.Sprint(product.Version)
	newVersion := fmt.Sprint(product.Version + 1)

	expressionAttributeValues := map[string]types.AttributeValue{
		":newQuantity": &types.AttributeValueMemberN{Value: newQuantity},
		":oldVersion":  &types.AttributeValueMemberN{Value: oldVersion},
		":newVersion":  &types.AttributeValueMemberN{Value: newVersion},
	}

	_, err := d.UpdateItem(tableName, key, updateExpression, conditionExpression, expressionAttributeValues)

	if err != nil {
		log.Printf("WARN: Couldn't update product %v\n", product.Id)
		return err
	}

	return nil
}
