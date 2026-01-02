package reaction

import (
	"context"

	domainreaction "hub/internal/domain/reaction"
	"hub/internal/domain/track"
)

// CheckReactionHandler handles the check reaction use case.
type CheckReactionHandler struct {
	reactionRepo domainreaction.Repository
}

// NewCheckReactionHandler creates a new CheckReactionHandler.
func NewCheckReactionHandler(reactionRepo domainreaction.Repository) *CheckReactionHandler {
	return &CheckReactionHandler{reactionRepo: reactionRepo}
}

// Handle executes the check reaction use case.
func (h *CheckReactionHandler) Handle(ctx context.Context, query CheckReactionQuery) (*CheckReactionResult, error) {
	userID, err := domainreaction.NewUserID(query.UserID)
	if err != nil {
		return nil, err
	}

	trackID, err := track.NewTrackID(query.TrackID)
	if err != nil {
		return nil, err
	}

	reaction, err := h.reactionRepo.FindByUserAndTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	if reaction == nil {
		return &CheckReactionResult{
			HasReacted: false,
		}, nil
	}

	return &CheckReactionResult{
		HasReacted: true,
		Reaction:   reaction.ReactionType().String(),
	}, nil
}
