package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"567_final/db"
	"567_final/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.ConnectMongoDB()

	r := mux.NewRouter()
	r.HandleFunc("/posts", handlers.GetPostsHandler).Methods("GET")
	r.HandleFunc("/log", handlers.GetLogHandler).Methods("GET")
	r.HandleFunc("/policy", handlers.GetPolicyHandler).Methods("GET")
	r.HandleFunc("/newPolicy", handlers.GetNewPolicyHandler).Methods("GET")
	r.HandleFunc("/simulation", handlers.GetSimulationHandler).Methods("GET")
	r.HandleFunc("/concern", handlers.PostConcernHandler).Methods("POST")
	r.HandleFunc("/updatePolicy", handlers.UpdatePolicyHandler).Methods("POST")
	r.HandleFunc("/updatePost", handlers.UpdatePostHandler).Methods("POST")
	r.HandleFunc("/log", handlers.DeleteAllLogsHandler).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
