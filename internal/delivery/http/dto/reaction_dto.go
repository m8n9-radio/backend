package dto

type ReactionRequest struct {
	UserID  string `json:"user_id" validate:"required,len=32"`
	TrackID string `json:"track_id" validate:"required,len=32"`
}

type ReactionResponse struct {
	Success bool `json:"success"`
}

type CheckReactionResponse struct {
	HasReacted bool   `json:"has_reacted"`
	Reaction   string `json:"reaction,omitempty"`
}
