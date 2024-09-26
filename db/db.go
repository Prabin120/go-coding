package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	QuestionsCollection *mongo.Collection
	TestCasesCollection *mongo.Collection
	client              *mongo.Client // Move the client to a package-level variable
)

// ConnectDB establishes a connection to MongoDB and returns a client.
func ConnectDB() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatal(err)
	}

	QuestionsCollection = client.Database("code_compiler").Collection("questions")
	TestCasesCollection = client.Database("code_compiler").Collection("testcases")
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

// DisconnectDB closes the MongoDB client connection.
func DisconnectDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
