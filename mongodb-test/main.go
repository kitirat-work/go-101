package main

import (
	"context"
	"fmt"
	"log"

	"your_project/mongotestcontainer"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// Start MongoDB Test Container
	mongoContainer, err := mongotestcontainer.StartMongoContainer(ctx)
	if err != nil {
		log.Fatalf("Error starting MongoDB container: %v", err)
	}
	defer mongoContainer.StopMongoContainer(ctx)

	fmt.Println("MongoDB running at:", mongoContainer.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoContainer.URI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Start a session for transactions
	session, err := client.StartSession()
	if err != nil {
		log.Fatalf("Error starting session: %v", err)
	}
	defer session.EndSession(ctx)

	// Perform a transaction
	err = session.StartTransaction()
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}

	collection := client.Database("test").Collection("users")

	_, err = collection.InsertOne(ctx, map[string]interface{}{"name": "Alice"})
	if err != nil {
		log.Fatalf("Error inserting document: %v", err)
	}

	err = session.CommitTransaction(ctx)
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	fmt.Println("Transaction committed successfully!")
}
