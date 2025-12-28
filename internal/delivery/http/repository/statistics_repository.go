package repository

import (
	"context"

	"hub/internal/delivery/http/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	StatisticsRepository interface {
		GetHistory(ctx context.Context) ([]*entity.Track, error)
		GetTopListened(ctx context.Context) ([]*entity.Track, error)
		GetTopRotate(ctx context.Context) ([]*entity.Track, error)
		GetTopLikes(ctx context.Context) ([]*entity.Track, error)
		GetTopDislikes(ctx context.Context) ([]*entity.Track, error)
	}
	statisticsRepository struct {
		pool *pgxpool.Pool
	}
)

func NewStatisticsRepository(pool *pgxpool.Pool) StatisticsRepository {
	return &statisticsRepository{pool: pool}
}

func (r *statisticsRepository) GetHistory(ctx context.Context) ([]*entity.Track, error) {
	query := `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks
		ORDER BY created_at DESC
		LIMIT 5
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*entity.Track
	for rows.Next() {
		var track entity.Track
		err := rows.Scan(
			&track.Title,
			&track.Cover,
			&track.Rotate,
			&track.Likes,
			&track.Dislikes,
			&track.Listeners,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *statisticsRepository) GetTopListened(ctx context.Context) ([]*entity.Track, error) {
	query := `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks
		WHERE listeners > 0
		ORDER BY listeners DESC
		LIMIT 5
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*entity.Track
	for rows.Next() {
		var track entity.Track
		err := rows.Scan(
			&track.Title,
			&track.Cover,
			&track.Rotate,
			&track.Likes,
			&track.Dislikes,
			&track.Listeners,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *statisticsRepository) GetTopRotate(ctx context.Context) ([]*entity.Track, error) {
	query := `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks
		ORDER BY rotate DESC
		LIMIT 5
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*entity.Track
	for rows.Next() {
		var track entity.Track
		err := rows.Scan(
			&track.Title,
			&track.Cover,
			&track.Rotate,
			&track.Likes,
			&track.Dislikes,
			&track.Listeners,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *statisticsRepository) GetTopLikes(ctx context.Context) ([]*entity.Track, error) {
	query := `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks
		WHERE likes > 0
		ORDER BY likes DESC
		LIMIT 5
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*entity.Track
	for rows.Next() {
		var track entity.Track
		err := rows.Scan(
			&track.Title,
			&track.Cover,
			&track.Rotate,
			&track.Likes,
			&track.Dislikes,
			&track.Listeners,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *statisticsRepository) GetTopDislikes(ctx context.Context) ([]*entity.Track, error) {
	query := `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks
		WHERE dislikes > 0
		ORDER BY dislikes DESC
		LIMIT 5
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*entity.Track
	for rows.Next() {
		var track entity.Track
		err := rows.Scan(
			&track.Title,
			&track.Cover,
			&track.Rotate,
			&track.Likes,
			&track.Dislikes,
			&track.Listeners,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}
