package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	CorrelationIDHeader = "X-Correlation-ID"
	CorrelationIDKey    = "correlation_id"
)

// CorrelationIDMiddleware adds correlation ID to requests.
func CorrelationIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		correlationID := c.Get(CorrelationIDHeader)

		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		c.Locals(CorrelationIDKey, correlationID)
		c.Set(CorrelationIDHeader, correlationID)

		return c.Next()
	}
}

// GetCorrelationID retrieves the correlation ID from context.
func GetCorrelationID(c *fiber.Ctx) string {
	if id, ok := c.Locals(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}
