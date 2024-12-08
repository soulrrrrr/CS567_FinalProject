package handlers

import (
	"567_final/db"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// postVote
type UpdatePolicyRequest struct {
	ID      primitive.ObjectID `json:"_id"`
	UserID  string             `json:"userID"`
	Vote    int                `json:"vote"`
	Comment string             `json:"comment"`
}

type VoteResponse struct {
	Success bool `json:"success"`
}

func UpdatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	var request UpdatePolicyRequest

	// Read the raw request body to use in logging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	// Restore the body for further processing
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: insert vote to database
	collection := db.GetCollection("uiuc-policy")

	// Check if the policy exists
	filter := bson.M{"_id": request.ID, "is_final": false}
	var policy bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "This is an exist policy or there is no this policy. You can only comment/vote to new policies.", http.StatusNotFound)
			return
		}
		http.Error(w, "Error finding policy", http.StatusInternalServerError)
		return
	}

	if policy["comments"] == nil {
		fmt.Println("Comments field is null, setting up as an array")

		// Define the update to set 'comments' as an empty array
		update := bson.M{
			"$set": bson.M{"comments": bson.A{}},
		}

		// Perform the update
		_, err = collection.UpdateOne(
			context.TODO(),
			filter,
			update,
		)
		if err != nil {
			fmt.Printf("Failed to initialize 'comments' as an array: %v\n", err)
			http.Error(w, "Failed to initialize 'comments'", http.StatusInternalServerError)
			return
		}
	}

	newComment := db.Comment{
		Author:    request.UserID,
		Body:      request.Comment,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}

	// Log the request and response
	var parsedRequest map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &parsedRequest); err != nil {
		fmt.Printf("Failed to parse request for logging: %v\n", err)
	}

	update := bson.M{
		"$inc": bson.M{"vote_count": request.Vote}, // Increment the vote count
	}

	// Add a new comment if provided
	if request.Comment != "" {
		update["$push"] = bson.M{
			"comments": newComment, // Push the new comment to the array
		}
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		filter,
		update,
	)

	if err != nil {
		if writeException, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range writeException.WriteErrors {
				fmt.Printf("Write error: %v\n", writeErr.Message)
			}
		}
		http.Error(w, "Failed to update the document", http.StatusInternalServerError)
		return
	}

	// Create a response based on the received data
	response := map[string]interface{}{
		"success": true,
		"message": "Policy updated successfully",
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	db.Logger.Log("UPDATEPOLICY", request.UserID, "Processed update policy request", parsedRequest, response)

}
