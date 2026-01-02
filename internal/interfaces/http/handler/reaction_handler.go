package handler

import (
	"errors"

	appreaction "hub/internal/application/reaction"
	"hub/internal/domain/reaction"
	"hub/internal/domain/track"
	"hub/internal/interfaces/http/dto"

	"github.com/gofiber/fiber/v2"
)

// ReactionHandler handles HTTP requests for reactions.
type ReactionHandler struct {
	addHandler   *appreaction.AddReactionHandler
	checkHandler *appreaction.CheckReactionHandler
}

// NewReactionHandler creates a new ReactionHandler.
func NewReactionHandler(
	addHandler *appreaction.AddReactionHandler,
	checkHandler *appreaction.CheckReactionHandler,
) *ReactionHandler {
	return &ReactionHandler{
		addHandler:   addHandler,
		checkHandler: checkHandler,
	}
}

// Like handles like requests.
func (h *ReactionHandler) Like(c *fiber.Ctx) error {
	return h.addReaction(c, "like")
}

// Dislike handles dislike requests.
func (h *ReactionHandler) Dislike(c *fiber.Ctx) error {
	return h.addReaction(c, "dislike")
}

// addReaction handles adding a reaction.
func (h *ReactionHandler) addReaction(c *fiber.Ctx, reactionType string) error {
	userID := c.Get("X-User-ID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("X-User-ID header is required"))
	}

	trackID := c.Params("trackId")
	if trackID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Track ID is required"))
	}

	err := h.addHandler.Handle(c.Context(), appreaction.AddReactionCommand{
		UserID:   userID,
		TrackID:  trackID,
		Reaction: reactionType,
	})

	if err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Check handles check reaction requests.
func (h *ReactionHandler) Check(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("X-User-ID header is required"))
	}

	trackID := c.Params("trackId")
	if trackID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Track ID is required"))
	}

	result, err := h.checkHandler.Handle(c.Context(), appreaction.CheckReactionQuery{
		UserID:  userID,
		TrackID: trackID,
	})

	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(dto.CheckReactionResponse{
		HasReacted: result.HasReacted,
		Reaction:   result.Reaction,
	})
}

// handleError maps domain errors to HTTP responses.
func (h *ReactionHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, reaction.ErrInvalidReactionType):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Invalid reaction type"))
	case errors.Is(err, reaction.ErrInvalidUserID):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Invalid user ID"))
	case errors.Is(err, reaction.ErrReactionExists):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrConflict("User has already reacted to this track"))
	case errors.Is(err, reaction.ErrTrackNotFound):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrNotFound("Track not found"))
	case errors.Is(err, track.ErrInvalidTrackID):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Invalid track ID format"))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternalServer)
	}
}
