package handlers

import (
	"567_final/db"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("policy")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var policyList []db.Policy
	var policy db.Policy
	for cursor.Next(ctx) {
		if err := cursor.Decode(&policy); err != nil {
			http.Error(w, "Failed to decode post", http.StatusInternalServerError)
			return
		}
		policyList = append(policyList, policy)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policyList)
}
