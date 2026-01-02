package track

import (
	"context"

	"hub/internal/domain/track"
)

// GetTrackHandler handles the get track use case.
type GetTrackHandler struct {
	repo track.Repository
}

// NewGetTrackHandler creates a new GetTrackHandler.
func NewGetTrackHandler(repo track.Repository) *GetTrackHandler {
	return &GetTrackHandler{repo: repo}
}

// Handle executes the get track use case.
func (h *GetTrackHandler) Handle(ctx context.Context, query GetTrackQuery) (*TrackDTO, error) {
	trackID, err := track.NewTrackID(query.ID)
	if err != nil {
		return nil, err
	}

	t, err := h.repo.FindByID(ctx, trackID)
	if err != nil {
		return nil, err
	}

	return &TrackDTO{
		ID:        t.ID().String(),
		Title:     t.Title().String(),
		Cover:     t.Cover().String(),
		Rotate:    t.Rotate(),
		Likes:     t.Likes(),
		Dislikes:  t.Dislikes(),
		Listeners: t.Listeners(),
	}, nil
}
