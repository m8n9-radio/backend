package handler

import (
	"errors"

	apptrack "hub/internal/application/track"
	"hub/internal/domain/track"
	"hub/internal/interfaces/http/dto"
	"hub/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// TrackHandler handles HTTP requests for tracks.
type TrackHandler struct {
	upsertHandler *apptrack.UpsertTrackHandler
	getHandler    *apptrack.GetTrackHandler
}

// NewTrackHandler creates a new TrackHandler.
func NewTrackHandler(
	upsertHandler *apptrack.UpsertTrackHandler,
	getHandler *apptrack.GetTrackHandler,
) *TrackHandler {
	return &TrackHandler{
		upsertHandler: upsertHandler,
		getHandler:    getHandler,
	}
}

// Upsert handles track creation/update requests.
func (h *TrackHandler) Upsert(c *fiber.Ctx) error {
	req, ok := c.Locals(middleware.TrackRequestKey).(*dto.CreateTrackRequest)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternalServer)
	}

	result, err := h.upsertHandler.Handle(c.Context(), apptrack.UpsertTrackCommand{
		ID:    req.Md5,
		Title: req.StreamTitle,
		Cover: req.StreamUrl,
	})

	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.TrackResponse{
		Rotate: result.Rotate,
	})
}

// Get handles get track requests.
func (h *TrackHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Track ID is required"))
	}

	result, err := h.getHandler.Handle(c.Context(), apptrack.GetTrackQuery{ID: id})
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(dto.GetTrackResponse{
		ID:        result.ID,
		Title:     result.Title,
		Cover:     result.Cover,
		Rotate:    result.Rotate,
		Likes:     result.Likes,
		Dislikes:  result.Dislikes,
		Listeners: result.Listeners,
	})
}

// handleError maps domain errors to HTTP responses.
func (h *TrackHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, track.ErrInvalidTrackID):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Invalid track ID format"))
	case errors.Is(err, track.ErrInvalidTitle):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrBadRequest("Invalid track title"))
	case errors.Is(err, track.ErrTrackNotFound):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrNotFound("Track not found"))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternalServer)
	}
}
