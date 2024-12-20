package routes

import (
	"code-compiler/internal/middlewares"
	"code-compiler/internal/usecases"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterQuestionRoutes(r *mux.Router, questionService *usecases.QuestionService) {
	wrappedCreateQuestion := middlewares.IsValidAdmin(http.HandlerFunc(questionService.CreateQuestion))
	wrappedGetQuestionById := middlewares.IsValidAdmin(http.HandlerFunc(questionService.GetQuestionById))
	wrappedUpdateQuestionById := middlewares.IsValidAdmin(http.HandlerFunc(questionService.UpdateQuestionById))
	wrappedCreateTestCases := middlewares.IsValidAdmin(http.HandlerFunc(questionService.CreateTestCase))
	wrappedUpdateTestCases := middlewares.IsValidAdmin(http.HandlerFunc(questionService.UpdateTestCases))
	wrappedGetQuestions := middlewares.IsValidUserWithoutReturn(http.HandlerFunc(questionService.GetQuestions))
	wrappedGetQuestionBySlug := middlewares.IsValidUserWithoutReturn(http.HandlerFunc(questionService.GetQuestionBySlug))
	r.Handle("/question", wrappedCreateQuestion).Methods(http.MethodPost)
	r.Handle("/question", wrappedGetQuestionById).Methods(http.MethodGet)
	r.Handle("/question", wrappedUpdateQuestionById).Methods(http.MethodPut)
	r.Handle("/questions", wrappedGetQuestions).Methods(http.MethodGet)
	r.Handle("/question/slug", wrappedGetQuestionBySlug).Methods(http.MethodGet)
	r.HandleFunc("/questions/tag", questionService.GetQuestionsByTag).Methods(http.MethodGet)
	r.Handle("/test-cases", wrappedCreateTestCases).Methods(http.MethodPost)
	r.HandleFunc("/test-cases", questionService.GetTestCases).Methods(http.MethodGet)
	r.Handle("/test-cases", wrappedUpdateTestCases).Methods(http.MethodPut)
	// r.HandleFunc("/validate-questions", questionService.GetInvalidQuestions).Methods(http.MethodPost)
}
