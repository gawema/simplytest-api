package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Collection
var mongoClient *mongo.Client

// ConnectDB initializes the connection to MongoDB
func ConnectDB() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB_NAME")
	collectionName := os.Getenv("MONGODB_COLLECTION")

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient = client
	db = client.Database(dbName).Collection(collectionName)
	log.Println("Connected to MongoDB!")
}

// CloseDB closes the MongoDB connection
func CloseDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Println("Error disconnecting from MongoDB:", err)
	} else {
		log.Println("Disconnected from MongoDB!")
	}
}
