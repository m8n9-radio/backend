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
		ExistsByMD5Sum(ctx context.Context, md5sum string) (bool, error)
		GetByMD5Sum(ctx context.Context, md5sum string) (*entity.Track, error)
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
		INSERT INTO tracks (id, md5sum, title, cover)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (md5sum) DO UPDATE SET
			cover = CASE
				WHEN (tracks.cover IS NULL OR tracks.cover = '') AND EXCLUDED.cover != ''
				THEN EXCLUDED.cover
				ELSE tracks.cover
			END,
			rotate = tracks.rotate + 1,
			updated_at = NOW()
		RETURNING id, md5sum, title, cover, rotate, likes, dislikes, created_at, updated_at
	`

	var result entity.Track
	err := r.pool.QueryRow(
		ctx,
		query,
		track.ID, track.MD5Sum, track.Title, track.Cover,
	).Scan(
		&result.ID,
		&result.MD5Sum,
		&result.Title,
		&result.Cover,
		&result.Rotate,
		&result.Likes,
		&result.Dislikes,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *trackRepository) ExistsByMD5Sum(ctx context.Context, md5sum string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tracks WHERE md5sum = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, md5sum).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r *trackRepository) GetByMD5Sum(ctx context.Context, md5sum string) (*entity.Track, error) {
	query := `
		SELECT id, md5sum, title, cover, rotate, likes, dislikes, created_at, updated_at
		FROM tracks WHERE md5sum = $1
	`

	var result entity.Track
	err := r.pool.QueryRow(ctx, query, md5sum).Scan(
		&result.ID,
		&result.MD5Sum,
		&result.Title,
		&result.Cover,
		&result.Rotate,
		&result.Likes,
		&result.Dislikes,
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
