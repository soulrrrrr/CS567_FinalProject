package handlers

import (
	"567_final/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func PostConcernHandler(w http.ResponseWriter, r *http.Request) {
	var request db.ConcernRequest

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Log the received data (optional)
	fmt.Printf("Received POST request: %+v\n", request)

	// Create a response based on the received data
	response := db.ConcernResponse{
		Policy: fmt.Sprintf("Policy generated for post %s with concern: %s", request.ID.Hex(), request.Concern),
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

}
