package handler

import (
	"errors"

	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/middleware"
	"hub/internal/delivery/http/service"

	"github.com/gofiber/fiber/v2"
)

type (
	ReactionHandler interface {
		Like(c *fiber.Ctx) error
		Dislike(c *fiber.Ctx) error
		Check(c *fiber.Ctx) error
	}
	reactionHandler struct {
		service service.ReactionService
	}
)

func NewReactionHandler(svc service.ReactionService) ReactionHandler {
	return &reactionHandler{service: svc}
}

func (h *reactionHandler) Like(c *fiber.Ctx) error {
	req, ok := c.Locals(middleware.ReactionRequestKey).(*dto.ReactionRequest)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get validated request",
		})
	}

	err := h.service.Like(c.Context(), req.UserID, req.TrackID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ReactionResponse{
		Success: true,
	})
}

func (h *reactionHandler) Dislike(c *fiber.Ctx) error {
	req, ok := c.Locals(middleware.ReactionRequestKey).(*dto.ReactionRequest)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get validated request",
		})
	}

	err := h.service.Dislike(c.Context(), req.UserID, req.TrackID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ReactionResponse{
		Success: true,
	})
}

func (h *reactionHandler) Check(c *fiber.Ctx) error {
	req, ok := c.Locals(middleware.ReactionRequestKey).(*dto.ReactionRequest)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get validated request",
		})
	}

	response, err := h.service.Check(c.Context(), req.UserID, req.TrackID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *reactionHandler) handleError(c *fiber.Ctx, err error) error {
	if errors.Is(err, service.ErrReactionAlreadyExists) {
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error:   "conflict",
			Message: "User has already reacted to this track",
		})
	}

	if errors.Is(err, service.ErrTrackNotFoundReaction) {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: "Track not found",
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
		Error:   "internal_error",
		Message: "Internal server error",
	})
}
