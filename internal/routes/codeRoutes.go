package routes

import (
	"code-compiler/internal/usecases"
	"net/http"
	"code-compiler/internal/middlewares"
	"github.com/gorilla/mux"
)

func RegisterCodeRoutes(r *mux.Router, codeRunService *usecases.CodeRunnerService) {
	wrappedRunTest := middlewares.IsValidUser(http.HandlerFunc(codeRunService.RunTest))
	wrappedSubmitTest := middlewares.IsValidUser(http.HandlerFunc(codeRunService.SubmitTest))
	r.Handle("/run-code", wrappedRunTest).Methods(http.MethodPost)
	r.Handle("/submit-code", wrappedSubmitTest).Methods(http.MethodPost)
}
