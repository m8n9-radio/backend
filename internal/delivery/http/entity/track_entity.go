package entity

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID        uuid.UUID
	MD5Sum    string
	Title     string
	Cover     string
	Rotate    int
	Likes     int
	Dislikes  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
