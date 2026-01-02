package dto

// ReactionRequest represents the HTTP request to add a reaction.
type ReactionRequest struct {
	UserID  string `json:"userId" validate:"required"`
	TrackID string `json:"trackId" validate:"required"`
}

// CheckReactionResponse represents the HTTP response for checking a reaction.
type CheckReactionResponse struct {
	HasReacted bool   `json:"hasReacted"`
	Reaction   string `json:"reaction,omitempty"`
}
