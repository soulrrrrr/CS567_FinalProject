package handlers

import "go.mongodb.org/mongo-driver/mongo"

var database *mongo.Database

func SetDatabase(db *mongo.Database) {
	database = db
}
