package main

import (
	"code-compiler/db"
	"code-compiler/repository"
	"code-compiler/usecases"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db.ConnectDB()

	r := mux.NewRouter()

	questionController := &repository.Question{}
	questionService := &usecases.QuestionService{Controller: questionController}

	r.HandleFunc("/", HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/question", questionService.CreateQuestion).Methods(http.MethodPost)
	r.HandleFunc("/question", questionService.GetQuestionById).Methods(http.MethodGet)
	r.HandleFunc("/question", questionService.UpdateQuestionById).Methods(http.MethodPut)
	r.HandleFunc("/questions", questionService.GetQuestions).Methods(http.MethodGet)
	r.HandleFunc("/questions/tag", questionService.GetQuestionsByTag).Methods(http.MethodGet)
	r.HandleFunc("/testcases", questionService.CreateTestCase).Methods(http.MethodPost)
	r.HandleFunc("/testcases", questionService.GetTestCases).Methods(http.MethodGet)

	// Set up CORS with the desired options
	corsOptions := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	fmt.Println("Start server on port 8080")
	srv := &http.Server{
		Addr:    ":8000",
		Handler: handlers.CORS(corsOptions, corsMethods, corsHeaders)(r),
	}

	// Channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error", err)
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
