package handlers

import (
	"567_final/db"
	"567_final/llmservice"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// postConcern
type ConcernRequest struct {
	PostID  primitive.ObjectID `json:"_id"`
	Concern string             `json:"concern"`
}

type ConcernResponse struct {
	Policy string `json:"policy"`
}

func PostConcernHandler(w http.ResponseWriter, r *http.Request) {
	var request ConcernRequest

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Log the received data (optional)
	fmt.Printf("Received POST request: %+v\n", request)

	// Run LLM to get new policy
	// Generate a new policy if needed
	policyList := GetPolicyFromDB(true)

	var currentPolicies []string
	for _, policy := range policyList {
		currentPolicies = append(currentPolicies, policy.PolicyDescription)
	}
	post, err := getPostFromDB(request.PostID)
	if err != nil {
		fmt.Printf("Error getting post: %v\n", err)
	}
	postContent := post.Body
	userComment := request.Concern

	fmt.Printf("Post Content: %s\n", postContent)

	newPolicy, err := llmservice.GenerateNewPolicy(currentPolicies, postContent, &userComment)
	if err != nil {
		fmt.Printf("Error generating new policy: %v\n", err)
		return
	}

	response := ConcernResponse{
		Policy: "",
	}

	if newPolicy != "" {
		fmt.Printf("New Policy Generated:\n%s\n\n", newPolicy)
		response.Policy = newPolicy

		insertPolicy := db.Policy{
			ID:                primitive.NewObjectID(),
			PolicyName:        newPolicy, // Assign policy name from newPolicy
			PolicyDescription: newPolicy, // You may want to improve this to a more detailed description
			VoteCount:         0,         // Default vote count
			IsFinal:           false,     // Default isFinal value
		}

		// Insert the new policy into the database
		err := addPolicyToDB(insertPolicy)
		if err != nil {
			fmt.Printf("Error inserting new policy: %v\n", err)
			return
		}

	} else {
		fmt.Println("No new policy needed based on the analysis.")
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

}

func getPostFromDB(id primitive.ObjectID) (db.Post, error) {
	collection := db.GetCollection("reddit-posts")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a filter to find the specific post
	filter := bson.M{"_id": id}

	var post db.Post
	err := collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		// Return zero-value post and error
		return db.Post{}, err
	}

	return post, nil
}

func addPolicyToDB(newPolicy db.Policy) error {
	collection := db.GetCollection("uiuc-policy")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the new policy into the collection
	_, err := collection.InsertOne(ctx, newPolicy)
	if err != nil {
		return err
	}

	return nil
}
