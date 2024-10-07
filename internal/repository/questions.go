package repository

import (
	"code-compiler/db"
	"code-compiler/internal/models"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Question struct {
	MongoCollection *mongo.Collection
	codeRunner      *CodeRunner
}

// CreateQuestion inserts a new question in the database.
func (r *Question) CreateQuestion(question *models.Question) (*models.Question, error) {
	question.ID = uuid.NewString() // Create a new ObjectID
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	_, err := db.QuestionsCollection.InsertOne(context.TODO(), question)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (r *Question) GetQuestionsByTag(tagName string) ([]models.Question, error) {
	var questions []models.Question
	cursor, err := db.QuestionsCollection.Find(context.TODO(), bson.M{"tags": bson.M{"$in": []string{tagName}}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &questions); err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *Question) GetQuestions() ([]models.Question, error) {
	var questions []models.Question
	cursor, err := db.QuestionsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &questions); err != nil {
		log.Fatal(err)
	}
	return questions, nil
}

func (r *Question) GetTestCases(questionId string) (*models.TestCase, error) {
	var testCases models.TestCase
	err := db.TestCasesCollection.FindOne(context.TODO(), bson.M{"questionId": questionId}).Decode(&testCases)
	if err != nil {
		return nil, err
	}
	return &testCases, nil
}

// GetQuestion fetches a question by its ID and populates its test cases.
func (r *Question) GetQuestionById(questionID string) (*models.Question, error) {
	var question models.Question
	err := db.QuestionsCollection.FindOne(context.TODO(), bson.M{"_id": questionID}).Decode(&question)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// CreateTestCase inserts a new test case in the database.
func (r *Question) CreateTestCase(testCase *models.TestCase) (*models.TestCase, error) {
	testCase.ID = uuid.NewString() // Create a new ObjectID
	testCase.CreatedAt = time.Now()
	testCase.UpdatedAt = time.Now()
	_, err := db.TestCasesCollection.InsertOne(context.TODO(), testCase)
	if err != nil {
		return nil, err
	}
	return testCase, nil
}

func (r *Question) UpdateQuestionById(questionID string, updatedData bson.M) (*models.Question, error) {
	updatedData["updatedAt"] = time.Now()
	err := db.QuestionsCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": questionID},
		bson.M{"$set": updatedData},
	).Decode(&updatedData)
	if err != nil {
		return nil, err
	}
	return r.GetQuestionById(questionID)
}

// CreateTestCase inserts a new test case in the database.
// func (r *Question) UpdateSolution(questionId string, code string) (*models.TestCase, error) {
// 	result, total, passed, err := r.codeRunner.ExecuteSubmit(commontypes.CodeRunnerType{
// 		Language: "go",
// 		Code: code,
// 		QuestionId: questionId,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if total != passed {
// 		return nil, fmt.Errorf("The code is failed in testCases. For input: ", result.Input)
// 	}

// 	return testCase, nil
// }
