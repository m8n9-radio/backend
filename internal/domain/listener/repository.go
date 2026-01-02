package listener

import "context"

// Repository defines the listener repository interface.
type Repository interface {
	// Save persists a listener.
	Save(ctx context.Context, listener *Listener) error

	// Exists checks if a listener exists for user and track.
	Exists(ctx context.Context, userID, trackID string) (bool, error)

	// CountByTrack returns the unique listener count for a track.
	CountByTrack(ctx context.Context, trackID string) (int, error)
}
