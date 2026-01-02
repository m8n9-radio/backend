package middleware

import (
	"strconv"
	"time"

	"hub/internal/infrastructure/metrics"

	"github.com/gofiber/fiber/v2"
)

// MetricsMiddleware records HTTP request metrics.
func MetricsMiddleware(m *metrics.Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		status := strconv.Itoa(c.Response().StatusCode())

		m.RecordHTTPRequest(c.Method(), c.Path(), status, duration)

		return err
	}
}
