package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InputOutput struct {
	Input  string `json:"input" bson:"input"`
	Output string `json:"expectedOutput" bson:"expectedOutput"`
}

// TestCase struct represents a single test case for a coding problem.
type TestCase struct {
	ID         string             `json:"_id,omitempty" bson:"_id"`
	QuestionID primitive.ObjectID `json:"questionID,omitempty" bson:"questionID"`
	IOPairs    []InputOutput      `json:"ioPairs,omitempty" bson:"ioPairs"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" bson:"updatedAt"`
}

//Have to fix the uniqueness and mandatory fields for the database.
