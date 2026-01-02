package reaction

import (
	"context"

	"hub/internal/domain/track"
)

// Repository defines the interface for reaction persistence.
type Repository interface {
	// Save persists a reaction.
	// Returns ErrReactionExists if the user has already reacted to the track.
	Save(ctx context.Context, reaction *Reaction) error

	// FindByUserAndTrack retrieves a reaction by user and track.
	// Returns nil if no reaction exists.
	FindByUserAndTrack(ctx context.Context, userID UserID, trackID track.TrackID) (*Reaction, error)

	// Exists checks if a reaction exists for the given user and track.
	Exists(ctx context.Context, userID UserID, trackID track.TrackID) (bool, error)
}
