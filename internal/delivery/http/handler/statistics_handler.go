package handler

import (
	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/service"

	"github.com/gofiber/fiber/v2"
)

type (
	StatisticsHandler interface {
		GetStatistics(c *fiber.Ctx) error
	}
	statisticsHandler struct {
		service service.StatisticsService
	}
)

func NewStatisticsHandler(svc service.StatisticsService) StatisticsHandler {
	return &statisticsHandler{service: svc}
}

func (h *statisticsHandler) GetStatistics(c *fiber.Ctx) error {
	statistics, err := h.service.GetStatistics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to fetch statistics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(statistics)
}
