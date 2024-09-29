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
	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	createdQuestion, err := svc.Controller.CreateQuestion(&question)
	if err != nil {
		http.Error(w, "Error creating question: "+err.Error(), http.StatusInternalServerError)
		return
	}
	res.Data = createdQuestion
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
func (svc *QuestionService) GetQuestionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}

	// Get the question ID from the URL parameters
	queryParams := r.URL.Query()
	questionID := queryParams.Get("id")
	if questionID == "" {
		http.Error(w, "question id required as query parameter", http.StatusBadRequest)
		return
	}
	// Call the controller to get the question by ID
	question, err := svc.Controller.GetQuestionById(questionID)
	if err != nil {
		http.Error(w, "Error fetching question: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Populate the response with the retrieved question
	res.Data = question

	// Send the question as the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (svc *QuestionService) GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}

	// Call the controller to get all questions
	questions, err := svc.Controller.GetQuestions()
	if err != nil {
		http.Error(w, "Error fetching questions: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Populate the response with the retrieved questions
	res.Data = questions
	// Send the questions as the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (svc *QuestionService) GetTestCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	vars := r.URL.Query()
	questionId := vars.Get("id")
	// Call the controller to get all questions
	tests, err := svc.Controller.GetTestCases(questionId)
	if err != nil {
		http.Error(w, "Error fetching questions: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Populate the response with the retrieved questions
	res.Data = tests
	// Send the questions as the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (svc *QuestionService) GetQuestionsByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	// Get the tag from the URL parameters
	vars := r.URL.Query()
	tag := vars.Get("tag")
	if tag == "" {
		http.Error(w, "tag string is required", http.StatusBadRequest)
		return
	}
	// Call the controller to get questions by the tag
	questions, err := svc.Controller.GetQuestionsByTag(tag)
	if err != nil {
		http.Error(w, "Error fetching questions by tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Populate the response with the retrieved questions
	res.Data = questions

	// Send the questions as the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (svc *QuestionService) UpdateQuestionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	// Get the question ID from the URL parameters
	vars := r.URL.Query()
	questionID := vars.Get("id")
	// Decode the incoming JSON request for the updated data
	var updatedData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the controller to update the question by ID
	updatedQuestion, err := svc.Controller.UpdateQuestionById(questionID, updatedData)
	if err != nil {
		http.Error(w, "Error updating question: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Populate the response with the updated question
	res.Data = updatedQuestion

	// Send the updated question as the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (svc *QuestionService) CreateTestCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}

	// Decode the incoming JSON request into the test case model
	var testCase models.TestCase
	if err := json.NewDecoder(r.Body).Decode(&testCase); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the controller to create the test case
	createdTestCase, err := svc.Controller.CreateTestCase(&testCase)
	if err != nil {
		http.Error(w, "Error creating test case: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Populate the response with the created test case
	res.Data = createdTestCase

	// Send the created test case as the response
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
