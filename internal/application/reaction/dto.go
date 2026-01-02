package reaction

// AddReactionCommand represents the command to add a reaction.
type AddReactionCommand struct {
	UserID   string
	TrackID  string
	Reaction string // "like" or "dislike"
}

// CheckReactionQuery represents the query to check a reaction.
type CheckReactionQuery struct {
	UserID  string
	TrackID string
}

// CheckReactionResult represents the result of checking a reaction.
type CheckReactionResult struct {
	HasReacted bool
	Reaction   string
}
