package routes

import (
	"code-compiler/internal/usecases"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterQuestionRoutes(r *mux.Router, questionService *usecases.QuestionService) {
	r.HandleFunc("/question", questionService.CreateQuestion).Methods(http.MethodPost)
	r.HandleFunc("/question", questionService.GetQuestionById).Methods(http.MethodGet)
	r.HandleFunc("/question", questionService.UpdateQuestionById).Methods(http.MethodPut)
	r.HandleFunc("/questions", questionService.GetQuestions).Methods(http.MethodGet)
	r.HandleFunc("/question/slug", questionService.GetQuestionBySlug).Methods(http.MethodGet)
	r.HandleFunc("/questions/tag", questionService.GetQuestionsByTag).Methods(http.MethodGet)
	r.HandleFunc("/test-cases", questionService.CreateTestCase).Methods(http.MethodPost)
	r.HandleFunc("/test-cases", questionService.GetTestCases).Methods(http.MethodGet)
	r.HandleFunc("/test-cases", questionService.UpdateTestCases).Methods(http.MethodPut)
}
