package handler

import (
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
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: fmt.Sprint(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(radioInfo)
}
