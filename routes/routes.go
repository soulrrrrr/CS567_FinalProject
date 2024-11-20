package routes

import (
	"567_final/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/posts/{index}", handlers.GetPostByIndexHandler).Methods("GET")
	return r
}
