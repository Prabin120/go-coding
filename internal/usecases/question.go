package usecases

import (
	"code-compiler/internal/models"
	"code-compiler/internal/repository"
	"encoding/json"
	"net/http"
)

type QuestionService struct {
	Controller *repository.Question
}

func (svc *QuestionService) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	createdQuestion, err := svc.Controller.CreateQuestion(&question)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	if res.Status {
		res.Data = createdQuestion
		res.Message = "Question added successfully"
		w.WriteHeader(http.StatusCreated)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (svc *QuestionService) GetQuestionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Get the question ID from the URL parameters
	queryParams := r.URL.Query()
	questionID := queryParams.Get("id")
	if questionID == "" {
		res.Status = false
		res.Message = "question doesn't found"
		w.WriteHeader(http.StatusBadRequest)
	}
	// Call the controller to get the question by ID
	question, err := svc.Controller.GetQuestionById(questionID)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Populate the response with the retrieved question
	if res.Status {
		res.Data = question
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *QuestionService) GetQuestionBySlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Get the question ID from the URL parameters
	queryParams := r.URL.Query()
	slug := queryParams.Get("slug")
	if slug == "" {
		res.Status = false
		res.Message = "slug required as query parameter"
		w.WriteHeader(http.StatusBadRequest)
	}
	// Call the controller to get the question by ID
	question, err := svc.Controller.GetQuestionBySlug(slug)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Populate the response with the retrieved question
	if res.Status {
		res.Data = question
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *QuestionService) GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Call the controller to get all questions
	questions, err := svc.Controller.GetQuestions()
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Populate the response with the retrieved questions
	if res.Status {
		res.Data = questions
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *QuestionService) GetTestCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	vars := r.URL.Query()
	questionId := vars.Get("questionId")
	// Call the controller to get all questions
	tests, errTest := svc.Controller.GetTestCases(questionId)
	question, errQuestion := svc.Controller.GetQuestionById(questionId)
	if errQuestion != nil {
		res.Status = false
		res.Message = errQuestion.Error()
		w.WriteHeader(http.StatusNotFound)
	}
	if errTest != nil {
		res.Message = errTest.Error()
		res.Data = question.Title
		w.WriteHeader(http.StatusNoContent)
	}
	// Populate the response with the retrieved questions
	if res.Status {
		res.Data = map[string]interface{}{
			"title": question.Title,
			"tests": tests,
		}
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *QuestionService) UpdateTestCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Get the question ID from the URL parameters
	vars := r.URL.Query()
	testCaseId := vars.Get("id")
	// Decode the incoming JSON request for the updated data
	var updatedData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Call the controller to update the question by ID
	updatedTestCase, err := svc.Controller.UpdateTestCases(testCaseId, updatedData)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Populate the response with the updated question
	if res.Status {
		res.Data = updatedTestCase
		w.WriteHeader(http.StatusCreated)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *QuestionService) GetQuestionsByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Get the tag from the URL parameters
	vars := r.URL.Query()
	tag := vars.Get("tag")
	if tag == "" {
		res.Status = false
		res.Message = "tag string is required"
		w.WriteHeader(http.StatusBadRequest)
	}
	// Call the controller to get questions by the tag
	questions, err := svc.Controller.GetQuestionsByTag(tag)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Populate the response with the retrieved questions
	if res.Status {
		res.Status = false
		res.Data = questions
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *QuestionService) UpdateQuestionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Get the question ID from the URL parameters
	vars := r.URL.Query()
	questionID := vars.Get("id")
	// Decode the incoming JSON request for the updated data
	var updatedData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Call the controller to update the question by ID
	updatedQuestion, err := svc.Controller.UpdateQuestionById(questionID, updatedData)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Populate the response with the updated question
	if res.Status {
		res.Data = updatedQuestion
		w.WriteHeader(http.StatusCreated)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *QuestionService) CreateTestCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Decode the incoming JSON request into the test case model
	var testCase models.TestCase
	if err := json.NewDecoder(r.Body).Decode(&testCase); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Call the controller to create the test case
	createdTestCase, err := svc.Controller.CreateTestCase(&testCase)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	if res.Status {
		res.Data = createdTestCase
		w.WriteHeader(http.StatusCreated)
	}
	// Send the created test case as the response
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}
