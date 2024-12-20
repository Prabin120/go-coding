package repository

import (
	"code-compiler/db"
	"code-compiler/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type Test struct {
	MongoCollection *mongo.Collection
}

func (r *Test) ValidateQuestions(questionIds []string) ([]string, error) {
	filter := bson.M{"_id": bson.M{"$in": questionIds}}
	cursor, err := db.QuestionsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	foundIDs := make(map[string]bool)
	for cursor.Next(context.TODO()) {
		var question models.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}
		foundIDs[question.ID] = true
	}
	var missingAptiIds []string
	for _, id := range questionIds {
		if !foundIDs[id] {
			missingAptiIds = append(missingAptiIds, id)
		}
	}
	return missingAptiIds, nil
}

func (r *Test) GetTestQuestions(questionIds []string) ([]models.Question, error) {
	filter := bson.M{"_id": bson.M{"$in": questionIds}}
	projection := bson.D{
		{"_id", 1},
		{"slug", 1},
		{"title", 1}, 
	}
	cursor, err := db.QuestionsCollection.Find(
		context.TODO(), 
		filter, 
		options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var questions []models.Question
	if err = cursor.All(context.TODO(), &questions); err != nil {
		return nil, err
	}
	return questions, nil
}