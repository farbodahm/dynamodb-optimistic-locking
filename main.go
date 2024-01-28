package main

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/ddb"
	"github.com/farbodahm/dynamodb-optimistic-locking/pkg/tables"
	"github.com/spf13/cobra"
)

// Make sure all of the purchase simulations are finished
var wg sync.WaitGroup

// numberOfFailedRequests counts number of failed requests because of ConditionalCheckFailedException error
var numberOfFailedRequests atomic.Int64

// getCommandLineParser creates the command line parser using Cobra
func getCommandLineParser() *cobra.Command {
	return &cobra.Command{
		Use:   "dynamodb-optimistic-locking",
		Short: "Simple application to create a race condition on DynamoDB and solve it using optimistic locking (versioning)",
	}
}

// simulateNewPurchase simulates a new purchase on DynamoDB using optimistic lock.
// Intentionally it waits for a second to make sure all of the requests get the same
// version from DDB to simulate a race condition on the table.
func simulateNewPurchase(d ddb.DynamoDB, tableName string, productId string, requestId int) {
	defer wg.Done()
	product, err := tables.GetProduct(d, tableName, productId)
	if err != nil {
		log.Fatalf("Failed to get product %v error: %v\n", productId, err)
	}

	// this sleep is used to intentionally create a race condition between different
	// goroutines trying to update the same product to test the optimistic lock mechanism.
	time.Sleep(2 * time.Second)

	product.Quantity -= 1
	if err = tables.SafeUpdateProductQuantity(d, tableName, product); err != nil {
		var e *types.ConditionalCheckFailedException
		if x := errors.As(err, &e); x {
			numberOfFailedRequests.Add(1)
		} else {
			log.Printf("WARN: Request Id %v failed with error %v", requestId, err)
		}
	}
}

func main() {
	dynamo := ddb.NewDynamoDBClient()
	cmd := getCommandLineParser()
	tableName := "products"

	var populateTable bool
	var numberOfRequests int
	cmd.Flags().BoolVar(&populateTable, "populate-table", false, "Populate the table with some sample data")
	cmd.Flags().IntVar(&numberOfRequests, "number-of-requests", 5, "Number of requests to simulate a concurrent access on DynamoDB")

	if err := cmd.Execute(); err != nil {
		log.Fatalln("Failed to parse arguments:", err)
	}

	if populateTable {
		log.Println("Populating `products` table with sample data...")
		if err := populateProductsTable(*dynamo); err != nil {
			log.Fatalln("Failed to populate products table:", err)
		}
	}

	wg.Add(numberOfRequests)
	for i := 0; i < numberOfRequests; i++ {
		go simulateNewPurchase(*dynamo, tableName, "p#1", i)
	}
	wg.Wait()

	log.Printf("Number of failed requests: %v\n", numberOfFailedRequests.Load())
}
