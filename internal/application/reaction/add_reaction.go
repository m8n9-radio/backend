package reaction

import (
	"context"

	appshared "hub/internal/application/shared"
	domainreaction "hub/internal/domain/reaction"
	"hub/internal/domain/track"
)

// AddReactionHandler handles the add reaction use case.
type AddReactionHandler struct {
	reactionRepo domainreaction.Repository
	trackRepo    track.Repository
	publisher    appshared.EventPublisher
}

// NewAddReactionHandler creates a new AddReactionHandler.
func NewAddReactionHandler(
	reactionRepo domainreaction.Repository,
	trackRepo track.Repository,
	publisher appshared.EventPublisher,
) *AddReactionHandler {
	return &AddReactionHandler{
		reactionRepo: reactionRepo,
		trackRepo:    trackRepo,
		publisher:    publisher,
	}
}

// Handle executes the add reaction use case.
func (h *AddReactionHandler) Handle(ctx context.Context, cmd AddReactionCommand) error {
	// Create value objects
	userID, err := domainreaction.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	trackID, err := track.NewTrackID(cmd.TrackID)
	if err != nil {
		return err
	}

	reactionType, err := domainreaction.NewReactionType(cmd.Reaction)
	if err != nil {
		return err
	}

	// Check if track exists
	exists, err := h.trackRepo.Exists(ctx, trackID)
	if err != nil {
		return err
	}
	if !exists {
		return domainreaction.ErrTrackNotFound
	}

	// Create reaction
	reaction := domainreaction.NewReaction(userID, trackID, reactionType)

	// Save reaction - handles duplicate check atomically via database constraint
	if err := h.reactionRepo.Save(ctx, reaction); err != nil {
		return err
	}

	// Publish event
	if h.publisher != nil {
		event := domainreaction.NewReactionAdded(userID, trackID, reactionType)
		if err := h.publisher.Publish(ctx, event); err != nil {
			// Log error but don't fail the operation
		}
	}

	return nil
}
