package postgres

import (
	"context"

	"hub/internal/application/listener"
)

// ListenerAdapter adapts Listentory for listener.Repository interface.
type ListenerAdapter struct {
	repo *ListenerRepository
}

// NewListenerAdapter creates a new adapter.
func NewListenerAdapter(repo *ListenerRepository) *ListenerAdapter {
	return &ListenerAdapter{repo: repo}
}

var _ listener.Repository = (*ListenerAdapter)(nil)

// TrackListener tracks a listener.
func (a *ListenerAdapter) TrackListener(ctx context.Context, userID, trackID string) error {
	return a.repo.TrackListener(ctx, userID, trackID)
}

// GetUniqueListenerCount returns unique listener count.
func (a *ListenerAdapter) GetUniqueListenerCount(ctx context.Context, trackID string) (int, error) {
	return a.repo.GetUniqueListenerCount(ctx, trackID)
}
