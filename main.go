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
	r.HandleFunc("/posts/{index}", handlers.GetPostByIndexHandler).Methods("GET")
	r.HandleFunc("/policy", handlers.GetPolicyHandler).Methods("GET")
	r.HandleFunc("/newPolicy", handlers.GetNewPolicyHandler).Methods("GET")
	r.HandleFunc("/concern", handlers.PostConcernHandler).Methods("POST")
	r.HandleFunc("/vote", handlers.PostVoteHandler).Methods("POST")
	r.HandleFunc("/simulation", handlers.PostSimulationHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
