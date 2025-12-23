package entity

import "time"

type ReactionType string

const (
	ReactionLike    ReactionType = "like"
	ReactionDislike ReactionType = "dislike"
)

type Reaction struct {
	ID        int64        `json:"id"`
	UserID    string       `json:"user_id"`
	TrackID   string       `json:"track_id"`
	Reaction  ReactionType `json:"reaction"`
	CreatedAt time.Time    `json:"created_at"`
}
