package postgres

import (
	"context"

	"hub/internal/application/statistics"

	"github.com/jackc/pgx/v5/pgxpool"
)

// StatisticsRepository implements statistics.Repository.
type StatisticsRepository struct {
	pool *pgxpool.Pool
}

// NewStatisticsRepository creates a new StatisticsRepository.
func NewStatisticsRepository(pool *pgxpool.Pool) *StatisticsRepository {
	return &StatisticsRepository{pool: pool}
}

var _ statistics.Repository = (*StatisticsRepository)(nil)

func (r *StatisticsRepository) GetHistory(ctx context.Context) ([]*statistics.TrackStats, error) {
	return r.queryTracks(ctx, `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks ORDER BY created_at DESC LIMIT 5
	`)
}

func (r *StatisticsRepository) GetTopListened(ctx context.Context) ([]*statistics.TrackStats, error) {
	return r.queryTracks(ctx, `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks WHERE listeners > 0 ORDER BY listeners DESC LIMIT 5
	`)
}

func (r *StatisticsRepository) GetTopRotate(ctx context.Context) ([]*statistics.TrackStats, error) {
	return r.queryTracks(ctx, `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks ORDER BY rotate DESC LIMIT 5
	`)
}

func (r *StatisticsRepository) GetTopLikes(ctx context.Context) ([]*statistics.TrackStats, error) {
	return r.queryTracks(ctx, `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks WHERE likes > 0 ORDER BY likes DESC LIMIT 5
	`)
}

func (r *StatisticsRepository) GetTopDislikes(ctx context.Context) ([]*statistics.TrackStats, error) {
	return r.queryTracks(ctx, `
		SELECT title, cover, rotate, likes, dislikes, listeners
		FROM tracks WHERE dislikes > 0 ORDER BY dislikes DESC LIMIT 5
	`)
}

func (r *StatisticsRepository) queryTracks(ctx context.Context, query string) ([]*statistics.TrackStats, error) {
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := make([]*statistics.TrackStats, 0)
	for rows.Next() {
		var t statistics.TrackStats
		if err := rows.Scan(&t.Title, &t.Cover, &t.Rotate, &t.Likes, &t.Dislikes, &t.Listeners); err != nil {
			return nil, err
		}
		tracks = append(tracks, &t)
	}
	return tracks, rows.Err()
}
