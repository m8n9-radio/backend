package repository

import (
	"context"
	"errors"

	"hub/internal/delivery/http/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrReactionAlreadyExists = errors.New("reaction already exists")
)

type (
	ReactionRepository interface {
		Create(ctx context.Context, reaction *entity.Reaction) error
		ExistsByUserAndTrack(ctx context.Context, userID, trackID string) (bool, error)
		GetByUserAndTrack(ctx context.Context, userID, trackID string) (*entity.Reaction, error)
	}
	reactionRepository struct {
		pool *pgxpool.Pool
	}
)

func NewReactionRepository(pool *pgxpool.Pool) ReactionRepository {
	return &reactionRepository{pool: pool}
}

func (r *reactionRepository) Create(ctx context.Context, reaction *entity.Reaction) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if reaction already exists
	exists, err := r.existsByUserAndTrackTx(ctx, tx, reaction.UserID, reaction.TrackID)
	if err != nil {
		return err
	}
	if exists {
		return ErrReactionAlreadyExists
	}

	// Insert reaction
	insertQuery := `
		INSERT INTO reactions (user_id, track_id, reaction)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, insertQuery, reaction.UserID, reaction.TrackID, reaction.Reaction)
	if err != nil {
		return err
	}

	// Update track counter based on reaction type
	var updateQuery string
	if reaction.Reaction == entity.ReactionLike {
		updateQuery = `UPDATE tracks SET likes = likes + 1, updated_at = NOW() WHERE id = $1`
	} else {
		updateQuery = `UPDATE tracks SET dislikes = dislikes + 1, updated_at = NOW() WHERE id = $1`
	}

	result, err := tx.Exec(ctx, updateQuery, reaction.TrackID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTrackNotFound
	}

	return tx.Commit(ctx)
}

func (r *reactionRepository) existsByUserAndTrackTx(ctx context.Context, tx pgx.Tx, userID, trackID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reactions WHERE user_id = $1 AND track_id = $2)`

	var exists bool
	err := tx.QueryRow(ctx, query, userID, trackID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *reactionRepository) ExistsByUserAndTrack(ctx context.Context, userID, trackID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reactions WHERE user_id = $1 AND track_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, trackID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *reactionRepository) GetByUserAndTrack(ctx context.Context, userID, trackID string) (*entity.Reaction, error) {
	query := `
		SELECT id, user_id, track_id, reaction, created_at
		FROM reactions
		WHERE user_id = $1 AND track_id = $2
	`

	var result entity.Reaction
	err := r.pool.QueryRow(ctx, query, userID, trackID).Scan(
		&result.ID,
		&result.UserID,
		&result.TrackID,
		&result.Reaction,
		&result.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
