package handler

import (
	"hub/internal/delivery/http/dto"
	fiberMiddleware "hub/internal/delivery/http/middleware"
	"hub/internal/delivery/http/service"

	"github.com/gofiber/fiber/v2"
)

type (
	TrackHandler interface {
		Upsert(c *fiber.Ctx) error
	}
	trackHandler struct {
		service service.TrackService
	}
)

func NewTrackHandler(svc service.TrackService) TrackHandler {
	return &trackHandler{service: svc}
}

func (h *trackHandler) Upsert(c *fiber.Ctx) error {
	req, ok := c.Locals(fiberMiddleware.TrackRequestKey).(*dto.CreateTrackRequest)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get validated request",
		})
	}

	track, err := h.service.Upsert(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Internal server error",
		})
	}

	response := dto.TrackResponse{
		ID:        track.ID.String(),
		MD5Sum:    track.MD5Sum,
		Title:     track.Title,
		Cover:     track.Cover,
		Rotate:    track.Rotate,
		Likes:     track.Likes,
		Dislikes:  track.Dislikes,
		CreatedAt: track.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: track.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
