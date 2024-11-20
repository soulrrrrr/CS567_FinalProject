package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Comment struct {
	Author    string `json:"author" bson:"author"`
	Body      string `json:"body" bson:"body"`
	CreatedAt string `json:"created_at" bson:"created_at"` // use string
}

type Post struct {
	ID        string    `json:"_id" bson:"_id"`
	Author    string    `json:"author" bson:"author"`
	Body      string    `json:"body" bson:"body"`
	Comments  []Comment `json:"comments" bson:"comments"`
	CreatedAt string    `json:"created_at" bson:"created_at"` // use string
	RedditID  string    `json:"id" bson:"id"`
	Permalink string    `json:"permalink" bson:"permalink"`
	Title     string    `json:"title" bson:"title"`
	Upvote    int       `json:"upvote" bson:"upvote"`
	URL       string    `json:"url" bson:"url"`
}

func GetPostByIndexHandler(w http.ResponseWriter, r *http.Request) {
	collection := database.Collection("reddit-posts")
	vars := mux.Vars(r)
	indexStr := vars["index"]
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var post Post
	i := 0
	for cursor.Next(ctx) {
		if err := cursor.Decode(&post); err != nil {
			http.Error(w, "Failed to decode post", http.StatusInternalServerError)
			return
		}
		if i == index {
			break
		}
		i++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
