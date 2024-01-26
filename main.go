package main

import (
	"log"

	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/ddb"
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/tables"
	"github.com/spf13/cobra"
)

// getCommandLineParser creates the command line parser using Cobra
func getCommandLineParser() *cobra.Command {
	return &cobra.Command{
		Use:   "dynamodb-optimistic-locking",
		Short: "Simple application to create a race condition on DynamoDB and solve it using optimistic locking (versioning)",
	}
}

func main() {
	dynamo := ddb.NewDynamoDBClient()
	cmd := getCommandLineParser()

	var populateTable bool
	cmd.Flags().BoolVar(&populateTable, "populate-table", false, "Populate the table with some sample data")

	if err := cmd.Execute(); err != nil {
		log.Fatalln("Failed to parse arguments:", err)
	}

	if populateTable {
		log.Println("Populating `products` table with sample data...")
		if err := populateProductsTable(*dynamo); err != nil {
			log.Fatalln("Failed to populate products table:", err)
		}
	}

	x, err := tables.GetProduct(*dynamo, "p#1")
	if err != nil {
		log.Fatalln("Failed to get product:", err)
	}
	log.Println(x)
}
