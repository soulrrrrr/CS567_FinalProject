package handlers

import (
	"567_final/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func PostSimulationHandler(w http.ResponseWriter, r *http.Request) {
	var request db.SimulationRequest

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Log the received data (optional)
	fmt.Printf("Received POST request: %+v\n", request)

	// TODO: connect with LLM and get feedback
	var results []db.SimulationResult

	for i := 0; i < 3; i++ {
		result := db.SimulationResult{
			Role:    fmt.Sprintf("role %d", i),
			Comment: fmt.Sprintf("I am role %d!", i),
		}
		results = append(results, result)
	}

	response := db.SimulationResponse{
		Results: results,
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

}
