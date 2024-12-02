package handlers

import (
	"567_final/db"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request body to get updated post data
	var updatedPost db.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure the provided _id exists in the request
	if updatedPost.ID.IsZero() {
		http.Error(w, "_id field is required", http.StatusBadRequest)
		return
	}

	// Update the post in the database
	collection := db.GetCollection("reddit-posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create filter for the post to update (match by _id)
	filter := bson.M{"_id": updatedPost.ID}

	// Create update data (you can choose which fields to update)
	update := bson.M{
		"$set": bson.M{
			"author":    updatedPost.Author,
			"body":      updatedPost.Body,
			"comments":  updatedPost.Comments,
			"upvote":    updatedPost.Upvote,
			"permalink": updatedPost.Permalink,
			"title":     updatedPost.Title,
			"url":       updatedPost.URL,
			"id":        updatedPost.RedditID,
		},
	}

	// Perform the update operation
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	response := map[string]bool{"success": true}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}
}
