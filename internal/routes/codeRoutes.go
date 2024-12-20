package routes

import (
	"code-compiler/internal/middlewares"
	"code-compiler/internal/usecases"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterCodeRoutes(r *mux.Router, codeRunService *usecases.CodeRunnerService) {
	wrappedRunTest := middlewares.IsValidUser(http.HandlerFunc(codeRunService.RunTest))
	wrappedSubmitTest := middlewares.IsValidUser(http.HandlerFunc(codeRunService.SubmitTest))
	wrappedGetSubmissions := middlewares.IsValidUser(http.HandlerFunc(codeRunService.GetUserSubmission))
	r.Handle("/run-code", wrappedRunTest).Methods(http.MethodPost)
	r.Handle("/code-submissions", wrappedGetSubmissions).Methods(http.MethodGet)
	r.Handle("/submit-code", wrappedSubmitTest).Methods(http.MethodPost)
}
