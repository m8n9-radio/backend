package postgres

import (
	"context"

	"hub/internal/application/listener"
)

// TrackListenerAdapter adapts TrackRepositstener.TrackRepository interface.
type TrackListenerAdapter struct {
	repo *TrackRepository
}

// NewTrackListenerAdapter creates a new adapter.
func NewTrackListenerAdapter(repo *TrackRepository) *TrackListenerAdapter {
	return &TrackListenerAdapter{repo: repo}
}

var _ listener.TrackRepository = (*TrackListenerAdapter)(nil)

// ExistsByID checks if a track exists by string ID.
func (a *TrackListenerAdapter) ExistsByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`
	var exists bool
	err := a.repo.pool.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

// UpdateListenerCount updates listener count by string ID.
func (a *TrackListenerAdapter) UpdateListenerCount(ctx context.Context, trackID string, count int) error {
	query := `UPDATE tracks SET listeners = $1, updated_at = NOW() WHERE id = $2`
	_, err := a.repo.pool.Exec(ctx, query, count, trackID)
	return err
}
