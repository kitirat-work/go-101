package mongotestcontainer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBContainer struct {
	Container testcontainers.Container
	URI       string
}

// StartMongoContainer starts a new MongoDB test container with a replica set.
func StartMongoContainer(ctx context.Context) (*MongoDBContainer, error) {
	fmt.Println("Starting MongoDB container...")
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		Cmd:          []string{"--replSet", "rs0", "--bind_ip_all"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	fmt.Printf("creates a generic container with parameters: %+v\n", req)
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	fmt.Println("get host where the container port is exposed")
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	fmt.Println("get mapped port")
	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%s/?directConnection=true", ip, port.Port())

	// Initialize Replica Set
	fmt.Printf("initialize replica set with URI: %s\n", uri)
	if err := initReplicaSet(ctx, uri); err != nil {
		return nil, fmt.Errorf("failed to initialize replica set: %w", err)
	}

	fmt.Println("MongoDB container started successfully!")

	return &MongoDBContainer{Container: container, URI: uri}, nil
}

// initReplicaSet initializes the MongoDB replica set.
func initReplicaSet(ctx context.Context, uri string) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(ctx)

	cmd := map[string]interface{}{
		"replSetInitiate": map[string]interface{}{
			"_id": "rs0",
			"members": []map[string]interface{}{
				{"_id": 0, "host": "localhost:27017"},
			},
		},
	}

	// Retry mechanism
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = client.Database("admin").RunCommand(ctx, cmd).Err()
		if err == nil {
			break
		}
		log.Printf("Retry %d/%d: failed to run replica set initiate command: %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to run replica set initiate command after %d retries: %w", maxRetries, err)
	}

	// Wait for replica set to initialize
	time.Sleep(5 * time.Second)
	return nil
}

// StopMongoContainer stops and removes the MongoDB test container.
func (m *MongoDBContainer) StopMongoContainer(ctx context.Context) {
	if err := m.Container.Terminate(ctx); err != nil {
		log.Printf("failed to terminate MongoDB container: %v", err)
	}
}
