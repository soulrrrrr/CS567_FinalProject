package handlers

import (
	"567_final/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func PostVoteHandler(w http.ResponseWriter, r *http.Request) {
	var request db.VoteRequest

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Log the received data (optional)
	fmt.Printf("Received POST request: %+v\n", request)

	// TODO: insert vote to database

	// Create a response based on the received data
	response := db.VoteResponse{
		// TODO: success depends on vote entered db or not
		Success: true,
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

}
