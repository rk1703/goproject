package main

import (
	"log"
	"net/http"

	"instagram-api/db"
	"instagram-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB("postgres://postgres:rk1703@localhost/instagram?sslmode=disable")

	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	r.HandleFunc("/posts/{id}", handlers.GetPost).Methods("GET")
	r.HandleFunc("/posts/users/{user_id}", handlers.ListUserPosts).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
