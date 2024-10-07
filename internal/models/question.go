package models

import (
	"time"
)

type Question struct {
	ID              string        `json:"_id,omitempty" bson:"_id"`
	Title           string        `json:"title,omitempty" bson:"title"`
	Description     string        `json:"description,omitempty" bson:"description"` // Difficulty level (easy, medium, hard)
	Difficulty      string        `json:"difficulty,omitempty" bson:"difficulty"`   // Tags for categorization (e.g., array, dynamic programming)
	Tags            []string      `json:"tags,omitempty" bson:"tags"`               // Reference to user who created the question
	SampleTestCases []InputOutput `json:"sampleTestCases,omitempty" bson:"sampleTestCases"`
	Solution        string        `json:"solution,omitempty" bson:"solution"`
	CreatedBy       string        `json:"createdBy,omitempty" bson:"createdBy"`             // Problem constraints (e.g., time complexity)
	TimeLimit       float64       `json:"timeLimit,omitempty" bson:"timeLimit"`             // Memory limit per test case execution (in kb)
	MemoryLimit     float64       `json:"memoryLimit,omitempty" bson:"memoryLimit"`         // Whether the question is public or private
	IsPublic        bool          `json:"isPublic,omitempty" bson:"isPublic"`               // Number of submissions for this question
	SubmissionCount int           `json:"submissionCount,omitempty" bson:"submissionCount"` // Success rate (in percentage)
	SuccessRate     float64       `json:"successRate,omitempty" bson:"successRate"`
	CreatedAt       time.Time     `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt,omitempty" bson:"updatedAt"`
}
