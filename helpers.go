package main

import (
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/ddb"
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/tables"
)

// populateProductsTable populates the products table with sample data.
func populateProductsTable(dynamo ddb.DynamoDB) error {

	sampleProducts := []tables.Products{
		{Id: "p#0", Name: "Product 0", Quantity: 10, Version: 1},
		{Id: "p#1", Name: "Product 1", Quantity: 4, Version: 1},
		{Id: "p#2", Name: "Product 2", Quantity: 2, Version: 1},
		{Id: "p#3", Name: "Product 3", Quantity: 6, Version: 1},
		{Id: "p#4", Name: "Product 4", Quantity: 1, Version: 1},
		{Id: "p#5", Name: "Product 5", Quantity: 0, Version: 1},
	}

	products := make([]ddb.DynamoDBWritable, len(sampleProducts))
	for i, product := range sampleProducts {
		products[i] = product
	}

	_, err := dynamo.AddBatch("products", products, 100)
	return err
}
