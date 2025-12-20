package repository

import (
	"context"
	"errors"

	"hub/internal/delivery/http/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrTrackNotFound = errors.New("track not found")
)

type (
	TrackRepository interface {
		Upsert(ctx context.Context, track *entity.Track) (*entity.Track, error)
		ExistsByID(ctx context.Context, id string) (bool, error)
		GetByID(ctx context.Context, id string) (*entity.Track, error)
	}
	trackRepository struct {
		pool *pgxpool.Pool
	}
)

func NewTrackRepository(pool *pgxpool.Pool) TrackRepository {
	return &trackRepository{pool: pool}
}

func (r *trackRepository) Upsert(ctx context.Context, track *entity.Track) (*entity.Track, error) {
	query := `
		INSERT INTO tracks (id, title, cover)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			cover = CASE
				WHEN (tracks.cover IS NULL OR tracks.cover = '') AND EXCLUDED.cover != ''
				THEN EXCLUDED.cover
				ELSE tracks.cover
			END,
			rotate = tracks.rotate + 1,
			updated_at = NOW()
		RETURNING id, title, cover, rotate, likes, dislikes, listeners, created_at, updated_at
	`

	var result entity.Track
	err := r.pool.QueryRow(
		ctx,
		query,
		track.ID, track.Title, track.Cover,
	).Scan(
		&result.ID,
		&result.Title,
		&result.Cover,
		&result.Rotate,
		&result.Likes,
		&result.Dislikes,
		&result.Listeners,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *trackRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r *trackRepository) GetByID(ctx context.Context, id string) (*entity.Track, error) {
	query := `
		SELECT id, title, cover, rotate, likes, dislikes, listeners, created_at, updated_at
		FROM tracks WHERE id = $1
	`

	var result entity.Track
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Title,
		&result.Cover,
		&result.Rotate,
		&result.Likes,
		&result.Dislikes,
		&result.Listeners,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTrackNotFound
		}
		return nil, err
	}

	return &result, nil
}
