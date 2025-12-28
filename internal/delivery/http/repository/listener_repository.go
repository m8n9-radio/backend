package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ListenerRepository interface {
		TrackListener(ctx context.Context, userID, trackID string) error
		GetUniqueListenerCount(ctx context.Context, trackID string) (int, error)
	}

	listenerRepository struct {
		pool *pgxpool.Pool
	}
)

func NewListenerRepository(pool *pgxpool.Pool) ListenerRepository {
	return &listenerRepository{pool: pool}
}

// TrackListener inserts a listener record if not exists (unique user-track combination)
func (r *listenerRepository) TrackListener(ctx context.Context, userID, trackID string) error {
	query := `
		INSERT INTO listeners (user_id, track_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, track_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, userID, trackID)
	return err
}

// GetUniqueListenerCount returns the number of unique listeners for a track
func (r *listenerRepository) GetUniqueListenerCount(ctx context.Context, trackID string) (int, error) {
	query := `SELECT COUNT(DISTINCT user_id) FROM listeners WHERE track_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, trackID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
