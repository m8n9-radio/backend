package middleware

import (
	"hub/internal/delivery/http/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const TrackRequestKey = "dto"

var validate = validator.New()

func ValidateCreateTrack() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.CreateTrackRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid request body",
			})
		}

		if err := validate.Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			message := formatValidationError(validationErrors[0])

			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error:   "validation_error",
				Message: message,
			})
		}

		c.Locals(TrackRequestKey, &req)
		return c.Next()
	}
}

func formatValidationError(fe validator.FieldError) string {
	switch fe.Field() {
	case "Md5":
		if fe.Tag() == "required" {
			return "Md5 is required"
		}
	case "StreamTitle":
		if fe.Tag() == "required" {
			return "StreamTitle is required"
		}
	}
	return "Validation failed for " + fe.Field()
}
