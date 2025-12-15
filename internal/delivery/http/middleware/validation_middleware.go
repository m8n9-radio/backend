package middleware

import (
	"hub/internal/delivery/http/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const TrackRequestKey = "validated_track_request"

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
	case "StreamTitle":
		switch fe.Tag() {
		case "required":
			return "StreamTitle is required"
		case "min":
			return "StreamTitle must be at least 1 character"
		case "max":
			return "StreamTitle must be at most 500 characters"
		}
	case "StreamUrl":
		switch fe.Tag() {
		case "required":
			return "StreamUrl is required"
		case "http_url":
			return "StreamUrl must be a valid HTTP URL"
		}
	case "MD5Sum":
		switch fe.Tag() {
		case "required":
			return "MD5Sum is required"
		case "len":
			return "MD5Sum must be 32 characters"
		case "hexadecimal":
			return "MD5Sum must be hexadecimal characters"
		}
	}
	return "Validation failed for " + fe.Field()
}
