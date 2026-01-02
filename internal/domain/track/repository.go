package track

import "context"

// Repository defines the interface for track persistence.
// This interface is defined in the domain layer and implemented in infrastructure.
type Repository interface {
	// Save persists a track aggregate.
	Save(ctx context.Context, track *Track) error

	// FindByID retrieves a track by its ID.
	// Returns ErrTrackNotFound if the track doesn't exist.
	FindByID(ctx context.Context, id TrackID) (*Track, error)

	// Exists checks if a track with the given ID exists.
	Exists(ctx context.Context, id TrackID) (bool, error)

	// UpdateListenerCount updates the listener count for a track.
	UpdateListenerCount(ctx context.Context, id TrackID, count int) error
}
