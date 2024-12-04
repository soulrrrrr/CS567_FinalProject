package handlers

import (
	"567_final/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// postVote
type UpdatePolicyRequest struct {
	ID      primitive.ObjectID `json:"_id"`
	UserID  string             `json:"user"`
	Vote    int                `json:"vote"`
	Comment string             `json:"comment"`
}

type VoteResponse struct {
	Success bool `json:"success"`
}

func UpdatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	var request UpdatePolicyRequest

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
			http.Error(w, "This is an exist policy. You can only comment/vote to new policies.", http.StatusNotFound)
			return
		}
		http.Error(w, "Error finding policy", http.StatusInternalServerError)
		return
	}

	newComment := db.Comment{
		Author:    request.UserID,
		Body:      request.Comment,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}
	// Add the vote (update or increment a field)
	update := bson.M{
		"$inc": bson.M{"vote_count": request.Vote}, // Increment a "votes" field (adjust as per your schema)
	}
	if request.Comment != "" {
		update["$push"] = bson.M{"comments": newComment}
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
