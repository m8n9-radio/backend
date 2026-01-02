package postgres

import (
	"context"
	"errors"
	"time"

	"hub/internal/domain/reaction"
	"hub/internal/domain/track"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReactionRepository implements reaction.Repository using PostgreSQL.
type ReactionRepository struct {
	pool *pgxpool.Pool
}

// NewReactionRepository creates a new ReactionRepository.
func NewReactionRepository(pool *pgxpool.Pool) *ReactionRepository {
	return &ReactionRepository{pool: pool}
}

// Save persists a reaction atomically with the track counter update.
func (r *ReactionRepository) Save(ctx context.Context, react *reaction.Reaction) error {
	// Start transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert with ON CONFLICT DO NOTHING to handle concurrent inserts gracefully
	query := `
		INSERT INTO reactions (user_id, track_id, reaction, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, track_id) DO NOTHING
	`

	result, err := tx.Exec(ctx, query,
		react.UserID().String(),
		react.TrackID().String(),
		react.ReactionType().String(),
		react.CreatedAt(),
	)

	if err != nil {
		return err
	}

	// Check if insert actually happened (rows affected = 0 means duplicate)
	if result.RowsAffected() == 0 {
		return reaction.ErrReactionExists
	}

	// Update track likes/dislikes count within same transaction
	var updateQuery string
	if react.IsLike() {
		updateQuery = `UPDATE tracks SET likes = likes + 1, updated_at = NOW() WHERE id = $1`
	} else {
		updateQuery = `UPDATE tracks SET dislikes = dislikes + 1, updated_at = NOW() WHERE id = $1`
	}

	_, err = tx.Exec(ctx, updateQuery, react.TrackID().String())
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(ctx)
}

// FindByUserAndTrack retrieves a reaction by user and track.
func (r *ReactionRepository) FindByUserAndTrack(ctx context.Context, userID reaction.UserID, trackID track.TrackID) (*reaction.Reaction, error) {
	query := `
		SELECT id, user_id, track_id, reaction, created_at
		FROM reactions WHERE user_id = $1 AND track_id = $2
	`

	var (
		id                               string
		userIDStr, trackIDStr, reactType string
		createdAt                        time.Time
	)

	err := r.pool.QueryRow(ctx, query, userID.String(), trackID.String()).
		Scan(&id, &userIDStr, &trackIDStr, &reactType, &createdAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return reaction.ReconstructReaction(id, userIDStr, trackIDStr, reactType, createdAt)
}

// Exists checks if a reaction exists for the given user and track.
func (r *ReactionRepository) Exists(ctx context.Context, userID reaction.UserID, trackID track.TrackID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reactions WHERE user_id = $1 AND track_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, userID.String(), trackID.String()).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
