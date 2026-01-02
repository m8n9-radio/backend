package reaction

import (
	"hub/internal/domain/shared"
	"hub/internal/domain/track"
)

// ReactionAdded is emitted when a user reacts to a track.
type ReactionAdded struct {
	shared.BaseEvent
	userID       UserID
	trackID      track.TrackID
	reactionType ReactionType
}

// NewReactionAdded creates a new ReactionAdded event.
func NewReactionAdded(userID UserID, trackID track.TrackID, reactionType ReactionType) ReactionAdded {
	return ReactionAdded{
		BaseEvent:    shared.NewBaseEvent("reaction.added"),
		userID:       userID,
		trackID:      trackID,
		reactionType: reactionType,
	}
}

// Payload returns the event data.
func (e ReactionAdded) Payload() interface{} {
	return map[string]interface{}{
		"user_id":       e.userID.String(),
		"track_id":      e.trackID.String(),
		"reaction_type": e.reactionType.String(),
	}
}

// UserID returns the user ID.
func (e ReactionAdded) UserID() UserID { return e.userID }

// TrackID returns the track ID.
func (e ReactionAdded) TrackID() track.TrackID { return e.trackID }

// ReactionType returns the reaction type.
func (e ReactionAdded) ReactionType() ReactionType { return e.reactionType }
