package usecases

import (
	commontypes "code-compiler/internal/commonTypes"
	"code-compiler/internal/middlewares"
	"code-compiler/internal/models"
	"code-compiler/internal/repository"
	"encoding/json"
	"net/http"
)

// CodeRunnerService struct to handle the business logic of code execution
type CodeRunnerService struct {
	Runner *repository.CodeRunner
}

type contextKey string

const userIDKey contextKey = "userID"

// interface conversion: interface {} is nil, not string

// RunTest handles incoming HTTP requests to run code and return test case results
func (svc *CodeRunnerService) RunTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	// Decode the incoming request body into CodeRunnerType struct
	var data commontypes.CodeRunnerType
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		res.Message = err.Error()
		res.Status = false
		w.WriteHeader(http.StatusBadRequest)
	}
	// Validate request data
	if data.Code == "" || data.Language == "" || data.QuestionId == "" {
		res.Message = "code, language, run type, and questionId are required"
		res.Status = false
		w.WriteHeader(http.StatusBadRequest)
	}
	// Call the repository function to execute the code
	result, err := svc.Runner.ExecuteTest(commontypes.CodeRunnerType{
		Language:   data.Language,
		Code:       data.Code,
		QuestionId: data.QuestionId,
	})

	if err != nil {
		res.Message = err.Error()
		res.Status = false
		w.WriteHeader(http.StatusInternalServerError)
	}
	if res.Status {
		res.Data = result
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Message = err.Error()
		res.Status = false
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *CodeRunnerService) SubmitTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	userId, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		res.Status = false
		res.Message = "User ID not found in context"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Status = true
	// Decode the incoming request body into CodeRunnerType struct
	var data commontypes.CodeRunnerType
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	}
	// Validate request data
	if data.Code == "" || data.Language == "" || data.QuestionId == "" {
		res.Status = false
		res.Message = "code, language, run type, and questionId are required"
		w.WriteHeader(http.StatusBadRequest)
	}
	// Call the repository function to execute the code
	failedCase, passedTestCases, totalTestCases, err := svc.Runner.ExecuteSubmit(commontypes.CodeRunnerType{
		UserId:     userId,
		Language:   data.Language,
		Code:       data.Code,
		QuestionId: data.QuestionId,
	})

	if err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	if res.Status {
		res.Data = map[string]interface{}{
			"failedCase":      failedCase,
			"passedTestCases": passedTestCases,
			"totalTestCases":  totalTestCases,
		}
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		res.Status = false
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *CodeRunnerService) GetUserSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	res.Status = true
	userId, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		res.Status = false
		res.Message = "User ID not found in context"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		return
	}
	// Get the tag from the URL parameters
	vars := r.URL.Query()
	question := vars.Get("question")
	if question == "" {
		res.Status = false
		res.Message = "question string is required"
		w.WriteHeader(http.StatusBadRequest)
	}
	// Call the controller to get questions by the tag
	questions, err := svc.Runner.GetUserSubmission(userId, question)
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
