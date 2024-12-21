package repository

import (
	"code-compiler/db"
	"code-compiler/internal/models"
	"code-compiler/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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
	slug, _ := r.GetQuestionBySlug(question.Slug, "")
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
func AddFieldsStage(userId string) bson.M {
	return bson.M{
		"$addFields": bson.M{
			"userStatus": bson.M{
				"$ifNull": bson.A{
					// Directly access the userId in the 'users' object, if it exists
					bson.M{
						"$ifNull": bson.A{
							"$users." + userId, // Access the user field by userId directly
							nil,                // Return null if the userId doesn't exist
						},
					},
					nil, // In case users field is missing or userId does not exist, return null
				},
			},
		},
	}
}

func (r *Question) GetQuestionBySlug(slug string, userId string) (*models.Question, error) {
	// Define the aggregation pipeline
	matchStage := bson.M{
		"$match": bson.M{
			"slug": slug, // Match question by slug
		},
	}
	addFieldsStage := AddFieldsStage(userId)
	projectStage := bson.M{
		"$project": bson.M{
			"users": 0, // Include userStatus in the projection
		},
	}
	if userId == "" {
		addFieldsStage = bson.M{}
	}
	pipeline := []bson.M{matchStage}
	if userId != "" {
		pipeline = append(pipeline, addFieldsStage)
	}
	pipeline = append(pipeline, projectStage)

	// Execute the aggregation query
	cursor, err := db.QuestionsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(context.TODO())

	// Store the results in a slice of questions
	var results []models.Question
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("failed to decode questions: %w", err)
	}

	// If no results are found, return nil
	if len(results) == 0 {
		return nil, nil // or any other error handling
	}

	// Return the first question (since you expect only one)
	return &results[0], nil
}

func (r *Question) GetQuestions(userId string, skip int, title string, difficulty string, status string) ([]models.Question, int64, error) {
	var questions []models.Question
	matchConditions := bson.M{}
	if title != "" {
		matchConditions["title"] = bson.M{"$regex": title, "$options": "i"} // Case-insensitive regex
	}
	if difficulty != "" {
		matchConditions["difficulty"] = difficulty
	}
	matchStage := bson.M{
		"$match": matchConditions,
	}
	addFieldsStage := AddFieldsStage(userId)

	pipeline := []bson.M{matchStage}
	if userId != "" {
		pipeline = append(pipeline, addFieldsStage)
	}

	if status != "" && userId != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"userStatus": status,
			},
		})
	}
	limit, err := strconv.Atoi(os.Getenv("PAGE_LIMIT"))
	if err != nil {
		limit = 10
	}
	paginationStages := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
	}
	projectStage := bson.M{
		"$project": bson.M{
			"_id":        1,
			"slug":       1,
			"title":      1,
			"difficulty": 1,
			"userStatus": 1,
		},
	}
	// pipeline = append(pipeline, projectStage)
	// Facet stage to calculate count and get paginated results
	pipeline = append(pipeline, bson.M{
		"$facet": bson.M{
			"metadata": []bson.M{
				{"$count": "total"}, // Count total documents
			},
			"data": append([]bson.M{
				projectStage,
			}, paginationStages...), // Apply projection and pagination
		},
	})
	// Use the Aggregate method instead of Find
	cursor, err := db.QuestionsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, 0, err // Return the error
	}
	defer cursor.Close(context.TODO())

	var results []struct {
		Metadata []struct {
			Total int64 `bson:"total"`
		} `bson:"metadata"`
		Data []models.Question `bson:"data"`
	}
	// Decode the cursor into the questions slice
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, 0, err // Return the error
	}
	var totalCount int64
	if len(results) > 0 && len(results[0].Metadata) > 0 {
		totalCount = results[0].Metadata[0].Total
		questions = results[0].Data
	} else {
		totalCount = 0
		questions = []models.Question{}
	}
	return questions, totalCount, nil
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
