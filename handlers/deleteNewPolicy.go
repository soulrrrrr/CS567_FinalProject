package handlers

// this is just for easy to test.

import (
	"567_final/db"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DeleteAllLogsHandler deletes all log entries from the database
func DeleteNewPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the logs collection
	collection := db.GetCollection("uiuc-policy")

	// Perform the delete operation
	result, err := collection.DeleteMany(ctx, bson.M{"is_final": false})
	if err != nil {
		http.Error(w, "Failed to delete policy", http.StatusInternalServerError)
		return
	}

	// Return the number of deleted documents
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"deleted": result.DeletedCount,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}
}
