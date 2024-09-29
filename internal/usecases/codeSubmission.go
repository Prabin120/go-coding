package usecases

import (
	commontypes "code-compiler/internal/commonTypes"
	"code-compiler/internal/models"
	"code-compiler/internal/repository"
	"encoding/json"
	"net/http"
)

// CodeRunnerService struct to handle the business logic of code execution
type CodeRunnerService struct {
	Runner *repository.CodeRunner
}

// RunTest handles incoming HTTP requests to run code and return test case results
func (svc *CodeRunnerService) RunTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	// Decode the incoming request body into CodeRunnerType struct
	var data commontypes.CodeRunnerType
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Validate request data
	if data.Code == "" || data.Language == "" || data.QuestionId == "" {
		http.Error(w, "code, language, run type, and questionId are required", http.StatusBadRequest)
		return
	}
	// Call the repository function to execute the code
	result, err := svc.Runner.ExecuteTest(commontypes.CodeRunnerType{
		Language:   data.Language,
		Code:       data.Code,
		QuestionId: data.QuestionId,
	})

	if err != nil {
		http.Error(w, "Error executing code: "+err.Error(), http.StatusInternalServerError)
		return
	}
	res.Data = result
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (svc *CodeRunnerService) SubmitTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &models.Response{}
	// Decode the incoming request body into CodeRunnerType struct
	var data commontypes.CodeRunnerType
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Validate request data
	if data.Code == "" || data.Language == "" || data.QuestionId == "" {
		http.Error(w, "code, language, run type, and questionId are required", http.StatusBadRequest)
		return
	}
	// Call the repository function to execute the code
	failedCase, passedTestCases, totalTestCases, err := svc.Runner.ExecuteSubmit(commontypes.CodeRunnerType{
		Language:   data.Language,
		Code:       data.Code,
		QuestionId: data.QuestionId,
	})

	if err != nil {
		http.Error(w, "Error executing code: "+err.Error(), http.StatusInternalServerError)
		return
	}
	res.Data = map[string]interface{}{
		"failedCase":      failedCase,
		"passedTestCases": passedTestCases,
		"totalTestCases":  totalTestCases,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
