package handlers

import (
	"567_final/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// postVote
type VoteRequest struct {
	ID     primitive.ObjectID `json:"_id"`
	UserID int                `json:"user"`
	Vote   int                `json:"vote"`
}

type VoteResponse struct {
	Success bool `json:"success"`
}

func PostVoteHandler(w http.ResponseWriter, r *http.Request) {
	var request VoteRequest

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Log the received data (optional)
	fmt.Printf("Received POST request: %+v\n", request)

	// TODO: insert vote to database
	collection := db.GetCollection("uiuc-policy")

	// Check if the policy exists
	filter := bson.M{"_id": request.ID, "is_final": false}
	var policy bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Policy not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error finding policy", http.StatusInternalServerError)
		return
	}

	// Add the vote (update or increment a field)
	update := bson.M{
		"$inc": bson.M{"vote_count": request.Vote}, // Increment a "votes" field (adjust as per your schema)
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update policy vote", http.StatusInternalServerError)
		return
	}

	// Create a response based on the received data
	response := VoteResponse{
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
