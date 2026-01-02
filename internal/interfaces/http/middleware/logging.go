package middleware

import (
	"time"

	"hub/internal/logger"

	"github.com/gofiber/fiber/v2"
)

// Logging creates a logging middleware.
func Logging(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		log.WithFields(map[string]interface{}{
			"method":   c.Method(),
			"path":     c.Path(),
			"status":   c.Response().StatusCode(),
			"duration": time.Since(start).String(),
			"ip":       c.IP(),
		}).Info("HTTP request")

		return err
	}
}
