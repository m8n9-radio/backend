package track

import (
	"context"
	"errors"

	appshared "hub/internal/application/shared"
	"hub/internal/domain/track"
)

// UpsertTrackHandler handles the upsert track use case.
type UpsertTrackHandler struct {
	repo      track.Repository
	publisher appshared.EventPublisher
}

// NewUpsertTrackHandler creates a new UpsertTrackHandler.
func NewUpsertTrackHandler(repo track.Repository, publisher appshared.EventPublisher) *UpsertTrackHandler {
	return &UpsertTrackHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle executes the upsert track use case.
func (h *UpsertTrackHandler) Handle(ctx context.Context, cmd UpsertTrackCommand) (*UpsertTrackResult, error) {
	// Create value objects
	trackID, err := track.NewTrackID(cmd.ID)
	if err != nil {
		return nil, err
	}

	title, err := track.NewTitle(cmd.Title)
	if err != nil {
		return nil, err
	}

	cover := track.NewCover(cmd.Cover)

	// Check if track exists
	existing, err := h.repo.FindByID(ctx, trackID)
	if err != nil && !errors.Is(err, track.ErrTrackNotFound) {
		return nil, err
	}

	var t *track.Track
	if existing != nil {
		// Update existing track
		t = existing
		t.IncrementRotation()
		t.UpdateCover(cover)
	} else {
		// Create new track
		t = track.NewTrack(trackID, title, cover)
	}

	// Persist
	if err := h.repo.Save(ctx, t); err != nil {
		return nil, err
	}

	// Publish events
	if h.publisher != nil && t.HasEvents() {
		if err := h.publisher.PublishAll(ctx, t.Events()); err != nil {
			// Log error but don't fail the operation
			// Events can be retried later
		}
		t.ClearEvents()
	}

	return &UpsertTrackResult{Rotate: t.Rotate()}, nil
}
