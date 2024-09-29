package models

import (
	"time"
)

// TestCase struct represents a single test case for a coding problem.
type TestCase struct {
	ID         string        `json:"_id,omitempty" bson:"_id"`
	QuestionID string        `json:"questionId,omitempty" bson:"questionId"`
	IOPairs    []InputOutput `json:"ioPairs,omitempty" bson:"ioPairs"`
	CreatedAt  time.Time     `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt,omitempty" bson:"updatedAt"`
}
