package handlers

import (
	"567_final/db"
	"567_final/logger"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GetAllLogs retrieves all logs from the database
func GetLogHandler(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("log")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		log.Printf("Error retrieving logs: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var logs []logger.LogEntry
	for cursor.Next(ctx) {
		var logEntry logger.LogEntry
		if err := cursor.Decode(&logEntry); err != nil {
			log.Printf("Error decoding log entry: %v\n", err)
			continue
		}
		logs = append(logs, logEntry)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error processing logs", http.StatusInternalServerError)
		log.Printf("Cursor error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		http.Error(w, "Failed to encode logs", http.StatusInternalServerError)
		log.Printf("Error encoding logs to JSON: %v\n", err)
	}
}
