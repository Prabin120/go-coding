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
	r.HandleFunc("/questions/tag", questionService.GetQuestionsByTag).Methods(http.MethodGet)
	r.HandleFunc("/testcases", questionService.CreateTestCase).Methods(http.MethodPost)
	r.HandleFunc("/testcases", questionService.GetTestCases).Methods(http.MethodGet)
}
