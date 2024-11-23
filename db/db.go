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

	createIndexes()

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func createIndexes() {
	testCaseIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "questionId", Value: 1}}, // Index on questionID
	}
	_, err := TestCasesCollection.Indexes().CreateOne(context.TODO(), testCaseIndexModel)
	if err != nil {
		log.Fatal("Failed to create index on TestCasesCollection: ", err)
	} else {
		fmt.Println("Unique index created on TestCasesCollection for questionID")
	}

	// Create unique index on title in QuestionsCollection
	questionTitleIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}}, // Index on title
		Options: options.Index().SetUnique(true), // Unique index
	}

	_, err = QuestionsCollection.Indexes().CreateOne(context.TODO(), questionTitleIndexModel)
	if err != nil {
		log.Fatal("Failed to create index on QuestionsCollection: ", err)
	} else {
		fmt.Println("Unique index created on QuestionsCollection for slug")
	}
}

// DisconnectDB closes the MongoDB client connection.
func DisconnectDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
