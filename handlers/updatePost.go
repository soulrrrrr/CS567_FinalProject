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
)

type UpdatePostRequest struct {
	UserID string  `json:"userID"`
	Post   db.Post `json:"post"`
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Read the raw request body to use in logging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	// Restore the body for further processing
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode the JSON request body
	var updatedPostRequest UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&updatedPostRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure the provided _id exists in the request
	if updatedPostRequest.Post.ID.IsZero() {
		http.Error(w, "_id field is required", http.StatusBadRequest)
		return
	}

	// Update the post in the database
	collection := db.GetCollection("reddit-posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create filter for the post to update (match by _id)
	filter := bson.M{"_id": updatedPostRequest.Post.ID}

	// Create update data
	update := bson.M{
		"$set": bson.M{
			"author":    updatedPostRequest.Post.Author,
			"body":      updatedPostRequest.Post.Body,
			"comments":  updatedPostRequest.Post.Comments,
			"upvote":    updatedPostRequest.Post.Upvote,
			"permalink": updatedPostRequest.Post.Permalink,
			"title":     updatedPostRequest.Post.Title,
			"url":       updatedPostRequest.Post.URL,
			"id":        updatedPostRequest.Post.RedditID,
		},
	}

	// Log the request and response
	var parsedRequest map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &parsedRequest); err != nil {
		fmt.Printf("Failed to parse request for logging: %v\n", err)
	}

	// Perform the update operation
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		db.Logger.Error(updatedPostRequest.UserID, err, parsedRequest, nil)
		return
	}

	// Prepare and return a success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"message": "Post updated successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		db.Logger.Error(updatedPostRequest.UserID, err, parsedRequest, nil)
		return
	}

	db.Logger.Log("UPDATEPOST", updatedPostRequest.UserID, "Processed update post request", parsedRequest, response)
}
