package main

import (
	"code-compiler/db"
	"code-compiler/internal/repository"
	"code-compiler/internal/routes"
	"code-compiler/internal/usecases"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db.ConnectDB()

	r := mux.NewRouter()

	questionController := &repository.Question{}
	questionService := &usecases.QuestionService{Controller: questionController}
	codeRunner := &repository.CodeRunner{}
	codeRunService := &usecases.CodeRunnerService{Runner: codeRunner}

	// Register routes from different files
	routes.RegisterQuestionRoutes(r, questionService)
	routes.RegisterCodeRoutes(r, codeRunService)

	r.HandleFunc("/", HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/test", runTest).Methods(http.MethodGet)

	// CORS configuration using rs/cors
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Allow credentials if needed
	})

	fmt.Println("Start server on port 8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler.Handler(r),
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()
	<-stop // Wait for an interrupt signal
	fmt.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Println("Server Shutdown:", err)
	}

	db.DisconnectDB() // Disconnect from MongoDB
	fmt.Println("Disconnected from MongoDB")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Its comming")
	fmt.Fprint(w, "Its working")
}

func runTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
}
