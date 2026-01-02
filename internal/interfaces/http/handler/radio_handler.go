package handler

import (
	"errors"

	"hub/internal/application/radio"
	"hub/internal/interfaces/http/dto"

	"github.com/gofiber/fiber/v2"
)

// RadioHandler handles HTTP requests for radio info.
type RadioHandler struct {
	service radio.Service
}

// NewRadioHandler creates a new RadioHandler.
func NewRadioHandler(svc radio.Service) *RadioHandler {
	return &RadioHandler{service: svc}
}

// GetInfo handles get radio info requests.
func (h *RadioHandler) GetInfo(c *fiber.Ctx) error {
	info, err := h.service.GetRadioInfo(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(dto.RadioResponse{
		Name:        info.Name,
		Description: info.Description,
		StreamUrl:   info.StreamUrl,
		Listener: dto.ListenerResponse{
			Current: info.Listeners,
			Peak:    info.ListenerPeak,
		},
	})
}

// GetListen handles get listener count requests.
func (h *RadioHandler) GetListen(c *fiber.Ctx) error {
	info, err := h.service.GetListeners(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(dto.ListenerResponse{
		Current: info.Current,
		Peak:    info.Peak,
	})
}

func (h *RadioHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, radio.ErrNoActiveStream):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrNotFound("No active stream available"))
	case errors.Is(err, radio.ErrInvalidResponse):
		return c.Status(fiber.StatusBadGateway).JSON(dto.NewErrorResponse("bad_gateway", "Invalid response from icecast server"))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternalServer)
	}
}
