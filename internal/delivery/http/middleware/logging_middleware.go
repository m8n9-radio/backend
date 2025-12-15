package middleware

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"hub/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		pid := os.Getpid()
		bodyBytes := c.Body()
		var requestBody interface{}
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
				requestBody = string(bodyBytes)
			}
		}

		queryParams := make(map[string]interface{})
		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
			keyStr := string(key)
			valueStr := string(value)

			if strings.HasPrefix(valueStr, "{") || strings.HasPrefix(valueStr, "[") {
				var jsonValue interface{}
				if err := json.Unmarshal([]byte(valueStr), &jsonValue); err == nil {
					queryParams[keyStr] = jsonValue
				} else {
					queryParams[keyStr] = valueStr
				}
			} else {
				queryParams[keyStr] = valueStr
			}
		})

		defer func() {
			latency := time.Since(start)

			var responseBody interface{}
			responseBytes := c.Response().Body()
			if len(responseBytes) > 0 {
				if jsonErr := json.Unmarshal(responseBytes, &responseBody); jsonErr != nil {
					responseBody = string(responseBytes)
				}
			}

			logEntry := log.WithFields(logrus.Fields{
				"latency":    latency.String(),
				"pid":        pid,
				"ip":         c.IP(),
				"user_agent": c.Get("User-Agent"),
				"method":     c.Method(),
				"url":        c.OriginalURL(),
				"body":       requestBody,
				"params":     queryParams,
				"response":   responseBody,
				"status":     c.Response().StatusCode(),
			})

			statusCode := c.Response().StatusCode()
			switch {
			case statusCode >= 500:
				logEntry.Error("HTTP Request - Server Error")
			case statusCode >= 400:
				logEntry.Warn("HTTP Request - Client Error")
			case statusCode >= 300:
				logEntry.Info("HTTP Request - Redirect")
			default:
				logEntry.Trace("HTTP Request - Success")
			}
		}()

		return c.Next()
	}
}
