package usecases

import (
	"code-compiler/internal/models"
	"code-compiler/internal/repository"
	"encoding/json"
	"net/http"
)

type TestService struct {
	Controller *repository.Test
}

// func (svc *QuestionService) CreateQuestion(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	res := &models.Response{}
// 	res.Status = true
// 	var question models.Question
// 	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
// 		res.Status = false
// 		res.Message = err.Error()
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	createdQuestion, err := svc.Controller.CreateQuestion(&question)
// 	if err != nil {
// 		res.Status = false
// 		res.Message = err.Error()
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	if res.Status {
// 		res.Data = createdQuestion
// 		res.Message = "Question added successfully"
// 		w.WriteHeader(http.StatusCreated)
// 	}
// 	if err := json.NewEncoder(w).Encode(res); err != nil {
// 		res.Status = false
// 		res.Message = err.Error()
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// }

func (svc *TestService) GetInvalidQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{
		Status: false,
	}
	var requestBody struct {
		Coding []string `json:"coding"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		res.Message = "Invalid request body: " + err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	if len(requestBody.Coding) == 0 {
		res.Message = "The 'coding' field is required and cannot be empty"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	invalidQuestions, err := svc.Controller.ValidateQuestions(requestBody.Coding)
	if err != nil {
		res.Message = "Failed to validate questions: " + err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	res.Data = invalidQuestions
	res.Status = len(invalidQuestions) == 0
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

type QuestionRequest struct {
    Questions []string `json:"questions"`
}

func (svc *TestService) GetTestQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	var req QuestionRequest
	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	if len(req.Questions) == 0 {
		res.Status = false
		res.Message = "Questions are not there"
		w.WriteHeader(http.StatusNoContent)
	}
	question, err := svc.Controller.GetTestQuestions(req.Questions)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
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