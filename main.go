package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

var collection *mongo.Collection

func connectMongoDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	// access the database URI from .env file
	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")
	collection = client.Database("project-cluster").Collection("reddit-posts")
}

// Get post by index
func getPostByIndexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse index from the URL
	vars := mux.Vars(r)
	indexStr := vars["index"]
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	// Query all posts and retrieve the i-th post
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// now I load 100 post from DB for each call, need to improve efficiency later
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
	fmt.Printf("[Query] %d\n", index)
	fmt.Println(post)
	fmt.Println("\n------")

	// Return the i-th post along with comments
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode post to JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	connectMongoDB()

	r := mux.NewRouter()
	// GET [index]th post
	r.HandleFunc("/posts/{index}", getPostByIndexHandler).Methods("GET")
	// POST user vote
	// POST user policy

	// use local database
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
