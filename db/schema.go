package db

import "go.mongodb.org/mongo-driver/bson/primitive"

// getNewPolicy
// getPolicy
type Policy struct {
	ID                string             `json:"_id" bson:"_id"`                               // MongoDB ObjectId, automatically assigned if not provided
	PolicyName        string             `json:"policy_name" bson:"policy_name"`               // Name of the policy
	PolicyDescription string             `json:"policy_description" bson:"policy_description"` // Description of the policy
	PostID            primitive.ObjectID `json:"post_id,omitempty" bson:"post_id,omitempty"`   // Post ID, converted from string to ObjectId
	VoteCount         int                `json:"vote_count" bson:"vote_count"`                 // Vote count for the policy
	IsFinal           bool               `json:"is_final" bson:"is_final"`                     // Indicates whether the policy is final
}

// getPostByIndex
type Comment struct {
	Author    string `json:"author" bson:"author"`
	Body      string `json:"body" bson:"body"`
	CreatedAt string `json:"created_at" bson:"created_at"` // use string
}

// getPostByIndex
type Post struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Author    string             `json:"author" bson:"author"`
	Body      string             `json:"body" bson:"body"`
	Comments  []Comment          `json:"comments" bson:"comments"`
	CreatedAt string             `json:"created_at" bson:"created_at"` // use string
	RedditID  string             `json:"id" bson:"id"`
	Permalink string             `json:"permalink" bson:"permalink"`
	Title     string             `json:"title" bson:"title"`
	Upvote    int                `json:"upvote" bson:"upvote"`
	URL       string             `json:"url" bson:"url"`
}
