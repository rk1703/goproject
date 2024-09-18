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

type PaginatedPostsResponse struct {
	Posts      []models.Post `json:"posts"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPosts int           `json:"total_posts"`
}

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

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // Default to the first page
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Default to 10 posts per page
	}

	offset := (page - 1) * limit

	// Query to get the total number of posts for the user
	var totalPosts int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = $1", userId).Scan(&totalPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query to get the posts with pagination
	rows, err := db.DB.Query("SELECT id, caption, image_url, posted_at FROM posts WHERE user_id = $1 LIMIT $2 OFFSET $3", userId, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Caption, &post.ImageURL, &post.PostedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	response := PaginatedPostsResponse{
		Posts:      posts,
		Page:       page,
		Limit:      limit,
		TotalPosts: totalPosts,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
