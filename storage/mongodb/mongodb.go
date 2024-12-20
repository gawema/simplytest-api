package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Represents a MongoDB connection with its client and collection
type MongoDB struct {
	client     *mongo.Client     // MongoDB client connection
	collection *mongo.Collection // Reference to the specific collection we're using
}

// Creates and initializes a new MongoDB connection
func NewMongoDB() (*MongoDB, error) {
	// Try to load .env file, but don't error if it doesn't exist
	_ = godotenv.Load()

	// Get MongoDB connection details from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("MONGODB_DB_NAME environment variable is not set")
	}

	collectionName := os.Getenv("MONGODB_COLLECTION")
	if collectionName == "" {
		return nil, fmt.Errorf("MONGODB_COLLECTION environment variable is not set")
	}

	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	// Set context with timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Get reference to collection
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDB{
		client:     client,
		collection: collection,
	}, nil
}

// Disconnects from MongoDB
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// Returns the MongoDB collection being used
func (m *MongoDB) Collection() *mongo.Collection {
	return m.collection
}
