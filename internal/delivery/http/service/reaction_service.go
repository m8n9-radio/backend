package service

import (
	"context"
	"errors"

	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/entity"
	"hub/internal/delivery/http/repository"
)

var (
	ErrReactionAlreadyExists = errors.New("user has already reacted to this track")
	ErrTrackNotFoundReaction = errors.New("track not found")
)

type (
	ReactionService interface {
		Like(ctx context.Context, userID, trackID string) error
		Dislike(ctx context.Context, userID, trackID string) error
		Check(ctx context.Context, userID, trackID string) (*dto.CheckReactionResponse, error)
	}
	reactionService struct {
		reactionRepo repository.ReactionRepository
		trackRepo    repository.TrackRepository
	}
)

func NewReactionService(reactionRepo repository.ReactionRepository, trackRepo repository.TrackRepository) ReactionService {
	return &reactionService{
		reactionRepo: reactionRepo,
		trackRepo:    trackRepo,
	}
}

func (s *reactionService) Like(ctx context.Context, userID, trackID string) error {
	return s.createReaction(ctx, userID, trackID, entity.ReactionLike)
}

func (s *reactionService) Dislike(ctx context.Context, userID, trackID string) error {
	return s.createReaction(ctx, userID, trackID, entity.ReactionDislike)
}

func (s *reactionService) createReaction(ctx context.Context, userID, trackID string, reactionType entity.ReactionType) error {
	// Check if track exists
	exists, err := s.trackRepo.ExistsByID(ctx, trackID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTrackNotFoundReaction
	}

	reaction := &entity.Reaction{
		UserID:   userID,
		TrackID:  trackID,
		Reaction: reactionType,
	}

	err = s.reactionRepo.Create(ctx, reaction)
	if err != nil {
		if errors.Is(err, repository.ErrReactionAlreadyExists) {
			return ErrReactionAlreadyExists
		}
		return err
	}

	return nil
}

func (s *reactionService) Check(ctx context.Context, userID, trackID string) (*dto.CheckReactionResponse, error) {
	reaction, err := s.reactionRepo.GetByUserAndTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	if reaction == nil {
		return &dto.CheckReactionResponse{
			HasReacted: false,
		}, nil
	}

	return &dto.CheckReactionResponse{
		HasReacted: true,
		Reaction:   string(reaction.Reaction),
	}, nil
}
