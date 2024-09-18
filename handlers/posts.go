package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"instagram-api/db"
	"instagram-api/models"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	json.NewDecoder(r.Body).Decode(&post)

	post.PostedAt = time.Now()

	_, err := db.DB.Exec("INSERT INTO posts (caption, image_url, posted_at, user_id) VALUES ($1, $2, $3, $4)",
		post.Caption, post.ImageURL, post.PostedAt, post.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post models.Post
	row := db.DB.QueryRow("SELECT id, caption, image_url, posted_at, user_id FROM posts WHERE id = $1", id)
	err := row.Scan(&post.ID, &post.Caption, &post.ImageURL, &post.PostedAt, &post.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func ListUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["user_id"])

	rows, err := db.DB.Query("SELECT id, caption, image_url, posted_at FROM posts WHERE user_id = $1 LIMIT 10", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.Caption, &post.ImageURL, &post.PostedAt)
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}
