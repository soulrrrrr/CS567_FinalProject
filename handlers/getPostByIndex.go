package handlers

import (
	"567_final/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPostByIndexHandler(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("reddit-posts")
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

	var post db.Post
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
	fmt.Printf("[GET] /posts %d\n", i)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
