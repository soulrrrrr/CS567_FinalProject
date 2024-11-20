package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"567_final/handlers"
	"567_final/routes"
	"567_final/utils"
)

func main() {
	database := utils.ConnectMongoDB()
	handlers.SetDatabase(database)

	r := routes.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
