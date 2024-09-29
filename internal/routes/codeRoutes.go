package routes

import (
	"code-compiler/internal/usecases"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterCodeRoutes(r *mux.Router, codeRunService *usecases.CodeRunnerService) {
	r.HandleFunc("/run-code", codeRunService.RunTest).Methods(http.MethodPost)
	r.HandleFunc("/submit-code", codeRunService.SubmitTest).Methods(http.MethodPost)
}
