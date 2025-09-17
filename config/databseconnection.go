package config

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	// Get Atlas URI from environment (copy it from Atlas → Connect → Drivers → Go)
	// Example: "mongodb+srv://<username>:<password>@<cluster-url>/?retryWrites=true&w=majority"
	uri := buildMongoURI()
	if uri == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable not set")
	}

	// Create client
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")
	return client, nil
}
func buildMongoURI() string {
	username := url.QueryEscape(os.Getenv("DB_USER"))
	password := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	cluster := os.Getenv("DB_CLUSTER")

	return fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
		username, password, cluster)
}
