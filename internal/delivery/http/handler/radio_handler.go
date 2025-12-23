package handler

import (
	"errors"

	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/service"

	"github.com/gofiber/fiber/v2"
)

type RadioHandler interface {
	GetInfo(c *fiber.Ctx) error
}

type radioHandler struct {
	service service.RadioService
}

func NewRadioHandler(svc service.RadioService) RadioHandler {
	return &radioHandler{service: svc}
}

func (h *radioHandler) GetInfo(c *fiber.Ctx) error {
	radioInfo, err := h.service.GetRadioInfo(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(radioInfo)
}

func (h *radioHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrIcecastUnavailable):
		return c.Status(fiber.StatusServiceUnavailable).JSON(dto.ErrorResponse{
			Error:   "service_unavailable",
			Message: "Icecast server unavailable",
		})
	case errors.Is(err, service.ErrInvalidResponse):
		return c.Status(fiber.StatusBadGateway).JSON(dto.ErrorResponse{
			Error:   "bad_gateway",
			Message: "Invalid response from Icecast",
		})
	case errors.Is(err, service.ErrNoActiveStream):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: "No active stream found",
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Internal server error",
		})
	}
}
