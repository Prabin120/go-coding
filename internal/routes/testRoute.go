package routes

import (
	"code-compiler/internal/usecases"
	"net/http"
	"code-compiler/internal/middlewares"
	"github.com/gorilla/mux"
)

func RegisterTestRoutes(r *mux.Router, testService *usecases.TestService) {
	wrappedValidateQuestions := middlewares.IsValidAdmin(http.HandlerFunc(testService.GetInvalidQuestions))
	r.HandleFunc("/test/questions", testService.GetTestQuestions).Methods(http.MethodPost)
	r.Handle("/test/validate-questions", wrappedValidateQuestions).Methods(http.MethodPost)
}
