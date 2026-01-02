package middleware

import (
	"hub/internal/interfaces/http/dto"

	"github.com/gofiber/fiber/v2"
)

// Context keys for validated data
const (
	TrackRequestKey = "track_request"
)

// Validator interface for request validation.
type Validator interface {
	Validate() error
}

// ValidateBody validates and parses request body.
func ValidateBody[T Validator]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "PARSE_ERROR",
					"message": "Invalid request body",
				},
			})
		}

		if err := body.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "VALIDATION_ERROR",
					"message": err.Error(),
				},
			})
		}

		c.Locals("body", &body)
		return c.Next()
	}
}

// GetBody retrieves the validated body from context.
func GetBody[T any](c *fiber.Ctx) (*T, bool) {
	body, ok := c.Locals("body").(*T)
	return body, ok
}

// ValidateTrackRequest validates track creation requests.
func ValidateTrackRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.CreateTrackRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if err := req.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Store in context using the DTO type
		c.Locals(TrackRequestKey, &req)

		return c.Next()
	}
}
