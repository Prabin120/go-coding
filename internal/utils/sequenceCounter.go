package utils

import (
	"context"
	"log"
	"strconv"

	"code-compiler/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Counter struct to store sequence values
type Counter struct {
	ID  string `bson:"_id"`
	Seq int64  `bson:"seq"`
}

// GetNextSequence fetches and increments the counter for the given collection
func GetNextSequence(collectionName string) (string, error) {
	filter := bson.M{"_id": collectionName}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counter Counter
	err := db.QuestionsCollection.Database().Collection("counters").FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&counter)
	if err != nil {
		log.Printf("Could not get next sequence: %v", err)
		return "", err
	}

	return strconv.FormatInt(counter.Seq, 10), nil
}
