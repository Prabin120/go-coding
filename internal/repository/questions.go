package repository

import (
	"code-compiler/db"
	"code-compiler/internal/models"
	"code-compiler/internal/utils"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Question struct {
	MongoCollection *mongo.Collection
}

// CreateQuestion inserts a new question in the database.
func (r *Question) CreateQuestion(question *models.Question) (*models.Question, error) {
	if question.Title == "" || question.Description == "" || question.Difficulty == "" ||
		question.MemoryLimit == 0.0 || question.Solution == "" || question.CodeTemplates == nil || question.SampleTestCases == nil || question.Tags == nil || question.TimeLimit == 0 {
		return nil, errors.New("please pass title, Description, Difficulty, MemoryLimit, Solution, CodeTemplate, SampleTestCases, Tags, TimeLimit")
	}
	slug, _ := r.GetQuestionBySlug(question.Slug)
	if slug != nil {
		return nil, errors.New("question already exists")
	}
	var seq, err_ = utils.GetNextSequence("question") // Create a new id
	if err_ != nil {
		return nil, errors.New("got error while creating id")
	}
	question.ID = seq
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()
	question.Slug = utils.MakeSlug(question.Title)
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

func (r *Question) GetQuestionBySlug(slug string) (*models.Question, error) {
	var question models.Question
	err := db.QuestionsCollection.FindOne(context.TODO(), bson.M{"slug": slug}).Decode(&question)
	if err != nil {
		return nil, err
	}
	return &question, nil
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

// GetQuestion fetches a question by its ID and populates its test cases.
func (r *Question) GetQuestionById(questionID string) (*models.Question, error) {
	var question models.Question
	err := db.QuestionsCollection.FindOne(context.TODO(), bson.M{"_id": questionID}).Decode(&question)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *Question) GetTestCases(questionId string) ([]models.InputOutput, error) {
	var testCases []struct {
		IOPairs []models.InputOutput
	}
	cursor, err := db.TestCasesCollection.Find(context.TODO(), bson.M{"questionId": questionId, "approved": true}, options.Find().SetProjection(bson.M{"ioPairs": 1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &testCases); err != nil {
		log.Fatal(err)
	}
	var ioPairs []models.InputOutput
	for _, testCase := range testCases {
		ioPairs = append(ioPairs, testCase.IOPairs...)
	}
	return ioPairs, nil
}

func (r *Question) GetTestCasesById(testCaseId string) (*models.TestCase, error) {
	var testCases models.TestCase
	err := db.TestCasesCollection.FindOne(context.TODO(), bson.M{"_id": testCaseId}).Decode(&testCases)
	if err != nil {
		return nil, err
	}
	return &testCases, nil
}

func (r *Question) UpdateTestCases(testCaseId string, updatedData bson.M) (*models.TestCase, error) {
	updatedData["updatedAt"] = time.Now()
	err := db.TestCasesCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": testCaseId},
		bson.M{"$set": updatedData},
	).Decode(&updatedData)
	if err != nil {
		return nil, err
	}
	return r.GetTestCasesById(testCaseId)
}

// CreateTestCase inserts a new test case in the database.
func (r *Question) CreateTestCase(testCase *models.TestCase) (*models.TestCase, error) {
	var seq, err = utils.GetNextSequence("testCase") // Create a new ObjectID
	if err != nil {
		return nil, errors.New("got error while creating id")
	}
	testCase.ID = seq
	testCase.CreatedAt = time.Now()
	testCase.UpdatedAt = time.Now()
	testCase.Approved = false
	_, err = db.TestCasesCollection.InsertOne(context.TODO(), testCase)
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
