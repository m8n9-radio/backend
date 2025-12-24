package handler

import (
	"errors"
	"fmt"
	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/service"

	"github.com/gofiber/fiber/v2"
)

type (
	RadioHandler interface {
		GetInfo(c *fiber.Ctx) error
	}
	radioHandler struct {
		service service.RadioService
	}
)

func NewRadioHandler(svc service.RadioService) RadioHandler {
	return &radioHandler{service: svc}
}

func (h *radioHandler) GetInfo(c *fiber.Ctx) error {
	radioInfo, err := h.service.GetRadioInfo(c.Context())
	if err != nil {
		if errors.Is(err, service.ErrNoActiveStream) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error:   "no_active_stream",
				Message: "No active stream available",
			})
		}
		if errors.Is(err, service.ErrInvalidResponse) {
			return c.Status(fiber.StatusBadGateway).JSON(dto.ErrorResponse{
				Error:   "invalid_response",
				Message: "Invalid response from icecast server",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: fmt.Sprint(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(radioInfo)
}
