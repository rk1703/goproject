package models

import "time"

type Post struct {
	ID       int       `json:"id"`
	Caption  string    `json:"caption"`
	ImageURL string    `json:"image_url"`
	PostedAt time.Time `json:"posted_at"`
	UserID   int       `json:"user_id"`
}
