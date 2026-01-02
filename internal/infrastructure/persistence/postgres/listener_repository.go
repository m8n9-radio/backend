package postgres

import (
	"context"

	"hub/internal/domain/listener"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ListenerRepository handles listener persistence.
type ListenerRepository struct {
	pool *pgxpool.Pool
}

// NewListenerRepository creates a new ListenerRepository.
func NewListenerRepository(pool *pgxpool.Pool) *ListenerRepository {
	return &ListenerRepository{pool: pool}
}

// Ensure ListenerRepository implements listener.Repository
var _ listener.Repository = (*ListenerRepository)(nil)

// Save persists a listener.
func (r *ListenerRepository) Save(ctx context.Context, l *listener.Listener) error {
	query := `
		INSERT INTO listeners (user_id, track_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, track_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, l.UserID(), l.TrackID(), l.CreatedAt())
	return err
}

// Exists checks if a listener exists.
func (r *ListenerRepository) Exists(ctx context.Context, userID, trackID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM listeners WHERE user_id = $1 AND track_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, trackID).Scan(&exists)
	return exists, err
}

// CountByTrack returns the unique listener count for a track.
func (r *ListenerRepository) CountByTrack(ctx context.Context, trackID string) (int, error) {
	query := `SELECT COUNT(DISTINCT user_id) FROM listeners WHERE track_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, trackID).Scan(&count)
	return count, err
}

// TrackListener tracks a listener for a track (legacy method for adapter).
func (r *ListenerRepository) TrackListener(ctx context.Context, userID, trackID string) error {
	l, err := listener.NewListener(userID, trackID)
	if err != nil {
		return err
	}
	return r.Save(ctx, l)
}

// GetUniqueListenerCount returns the number of unique listeners for a track (legacy).
func (r *ListenerRepository) GetUniqueListenerCount(ctx context.Context, trackID string) (int, error) {
	return r.CountByTrack(ctx, trackID)
}
