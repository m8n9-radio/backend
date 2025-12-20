package entity

import "time"

type Track struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Listeners int       `json:"listeners"`
	Rotate    int       `json:"rotate"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
