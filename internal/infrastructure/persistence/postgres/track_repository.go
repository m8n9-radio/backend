package postgres

import (
	"context"
	"errors"
	"time"

	"hub/internal/domain/track"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TrackRepository implements track.Repository using PostgreSQL.
type TrackRepository struct {
	pool *pgxpool.Pool
}

// NewTrackRepository creates a new TrackRepository.
func NewTrackRepository(pool *pgxpool.Pool) *TrackRepository {
	return &TrackRepository{pool: pool}
}

// Save persists a track aggregate.
func (r *TrackRepository) Save(ctx context.Context, t *track.Track) error {
	query := `
		INSERT INTO tracks (id, title, cover, rotate, likes, dislikes, listeners, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			cover = CASE
				WHEN (tracks.cover IS NULL OR tracks.cover = '') AND EXCLUDED.cover != ''
				THEN EXCLUDED.cover
				ELSE tracks.cover
			END,
			rotate = EXCLUDED.rotate,
			likes = EXCLUDED.likes,
			dislikes = EXCLUDED.dislikes,
			listeners = EXCLUDED.listeners,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(ctx, query,
		t.ID().String(),
		t.Title().String(),
		t.Cover().String(),
		t.Rotate(),
		t.Likes(),
		t.Dislikes(),
		t.Listeners(),
		t.CreatedAt(),
		time.Now(),
	)

	return err
}

// FindByID retrieves a track by its ID.
func (r *TrackRepository) FindByID(ctx context.Context, id track.TrackID) (*track.Track, error) {
	query := `
		SELECT id, title, cover, rotate, likes, dislikes, listeners, created_at, updated_at
		FROM tracks WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id.String())
	return r.scanTrack(row)
}

// Exists checks if a track with the given ID exists.
func (r *TrackRepository) Exists(ctx context.Context, id track.TrackID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id.String()).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UpdateListenerCount updates the listener count for a track.
func (r *TrackRepository) UpdateListenerCount(ctx context.Context, id track.TrackID, count int) error {
	query := `UPDATE tracks SET listeners = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, count, id.String())
	return err
}

// scanTrack scans a row into a Track aggregate.
func (r *TrackRepository) scanTrack(row pgx.Row) (*track.Track, error) {
	var (
		id, title, cover                   string
		rotate, likes, dislikes, listeners int
		createdAt, updatedAt               time.Time
	)

	err := row.Scan(&id, &title, &cover, &rotate, &likes, &dislikes, &listeners, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, track.ErrTrackNotFound
		}
		return nil, err
	}

	return track.ReconstructTrack(id, title, cover, rotate, likes, dislikes, listeners, createdAt, updatedAt)
}
