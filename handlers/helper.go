package handlers

import (
	"567_final/db"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// reuse functions for handlers
func GetPolicyFromDB(isFinal bool) []db.Policy {
	collection := db.GetCollection("uiuc-policy")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"is_final": isFinal})
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	var policyList []db.Policy
	var policy db.Policy
	for cursor.Next(ctx) {
		if err := cursor.Decode(&policy); err != nil {
			return nil
		}
		policyList = append(policyList, policy)
	}

	return policyList
}
