package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDBTransactions(t *testing.T) {
	ctx := context.Background()

	// Start the MongoDB container
	mongoC, mongoURI, err := setupMongoContainer()
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}
	defer mongoC.Terminate(ctx)

	// Wait for MongoDB to be ready
	time.Sleep(5 * time.Second)

	// Connect to MongoDB
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://test:testpassword@%s/?authSource=admin", mongoURI))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Get database and collection
	db := client.Database("testdb")
	accounts := db.Collection("accounts")

	// Insert initial documents
	_, err = accounts.InsertMany(ctx, []interface{}{
		bson.M{"_id": 1, "accountId": 1, "balance": 500},
		bson.M{"_id": 2, "accountId": 2, "balance": 500},
	})
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Start a transaction
	session, err := client.StartSession()
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	// Define the transaction function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Deduct money from one account
		_, err := accounts.UpdateOne(sessCtx, bson.M{"accountId": 1}, bson.M{"$inc": bson.M{"balance": -100}})
		if err != nil {
			return nil, err
		}

		// Add money to another account
		_, err = accounts.UpdateOne(sessCtx, bson.M{"accountId": 2}, bson.M{"$inc": bson.M{"balance": 100}})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Execute the transaction
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	// Verify the balances after transaction
	var result bson.M
	err = accounts.FindOne(ctx, bson.M{"accountId": 1}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to fetch account 1: %v", err)
	}
	if result["balance"].(int32) != 400 {
		t.Fatalf("Unexpected balance for account 1: %v", result["balance"])
	}

	err = accounts.FindOne(ctx, bson.M{"accountId": 2}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to fetch account 2: %v", err)
	}
	if result["balance"].(int32) != 600 {
		t.Fatalf("Unexpected balance for account 2: %v", result["balance"])
	}

	t.Log("Transaction test passed successfully!")
}

// Start MongoDB container with replica set
func setupMongoContainer() (testcontainers.Container, string, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0", // Use MongoDB 6.0
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "test",
			"MONGO_INITDB_ROOT_PASSWORD": "testpassword",
			"MONGO_REPLICA_SET_NAME":     "rs0",
		},
		Cmd: []string{"--replSet", "rs0"},
		WaitingFor: wait.ForLog("Waiting for connections").
			WithStartupTimeout(10 * time.Second),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	// Get container host & port
	endpoint, err := mongoC.Endpoint(ctx, "")
	if err != nil {
		return nil, "", err
	}

	// Initialize the Replica Set
	go func() {
		time.Sleep(3 * time.Second) // Wait for MongoDB to be ready
		clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://test:testpassword@%s", endpoint))
		client, err := mongo.Connect(ctx, clientOpts)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		defer client.Disconnect(ctx)

		adminDB := client.Database("admin")
		adminDB.RunCommand(ctx, bson.D{{Key: "replSetInitiate", Value: bson.D{}}})
	}()

	return mongoC, endpoint, nil
}
