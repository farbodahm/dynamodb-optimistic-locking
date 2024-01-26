package main

import (
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/ddb"
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/tables"
)

// populateProductsTable populates the Product table with sample data.
func populateProductsTable(dynamo ddb.DynamoDB) error {

	sampleProduct := []tables.Product{
		{Id: "p#0", Name: "Product 0", Quantity: 10, Version: 1},
		{Id: "p#1", Name: "Product 1", Quantity: 4, Version: 1},
		{Id: "p#2", Name: "Product 2", Quantity: 2, Version: 1},
		{Id: "p#3", Name: "Product 3", Quantity: 6, Version: 1},
		{Id: "p#4", Name: "Product 4", Quantity: 1, Version: 1},
		{Id: "p#5", Name: "Product 5", Quantity: 0, Version: 1},
	}

	Product := make([]ddb.DynamoDBWritable, len(sampleProduct))
	for i, product := range sampleProduct {
		Product[i] = product
	}

	_, err := dynamo.AddBatch("products", Product, 100)
	return err
}
