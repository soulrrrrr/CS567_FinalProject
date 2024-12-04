package logger

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type LogEntry struct {
	UserID    string                 `json:"userID"`
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Request   map[string]interface{} `json:"request"`
	Response  map[string]interface{} `json:"response"`
}

type MongoLogger struct {
	collection *mongo.Collection
}

func NewMongoLogger(client *mongo.Client, dbName, collectionName string) *MongoLogger {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoLogger{collection: collection}
}

func (l *MongoLogger) Log(level, userID, message string, request, response map[string]interface{}) {
	logEntry := LogEntry{
		UserID:    userID,
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Request:   request,
		Response:  response,
	}

	go func() { // Async logging to avoid blocking
		_, err := l.collection.InsertOne(context.Background(), logEntry)
		if err != nil {
			fmt.Printf("Failed to log entry: %v\n", err)
		}
	}()
}

func (l *MongoLogger) Error(userID string, err error, request, response map[string]interface{}) {
	l.Log("ERROR", userID, err.Error(), request, response)
}
