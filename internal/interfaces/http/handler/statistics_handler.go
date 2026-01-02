package handler

import (
	"hub/internal/application/statistics"
	"hub/internal/interfaces/http/dto"

	"github.com/gofiber/fiber/v2"
)

// StatisticsHandler handles HTTP requests for statistics.
type StatisticsHandler struct {
	service statistics.Service
}

// NewStatisticsHandler creates a new StatisticsHandler.
func NewStatisticsHandler(svc statistics.Service) *StatisticsHandler {
	return &StatisticsHandler{service: svc}
}

// GetStatistics handles get statistics requests.
func (h *StatisticsHandler) GetStatistics(c *fiber.Ctx) error {
	stats, err := h.service.GetStatistics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternalServer)
	}

	response := make([]*dto.StatisticCategory, len(stats))

	for i, cat := range stats {
		tracks := make([]*dto.TrackStats, len(cat.Tracks))
		for j, t := range cat.Tracks {
			tracks[j] = &dto.TrackStats{
				Title:     t.Title,
				Cover:     t.Cover,
				Rotate:    t.Rotate,
				Likes:     t.Likes,
				Dislikes:  t.Dislikes,
				Listeners: t.Listeners,
			}
		}
		response[i] = &dto.StatisticCategory{
			Key:    cat.Key,
			Icon:   cat.Icon,
			Tracks: tracks,
		}
	}

	return c.JSON(response)
}
