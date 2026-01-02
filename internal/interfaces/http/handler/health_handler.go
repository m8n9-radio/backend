package handler

import (
	"context"
	"fmt"
	"time"

	"hub/internal/interfaces/http/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *pgxpool.Pool, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

// Health handles GET /health requests.
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]dto.Check)
	allHealthy := true

	dbCheck := h.checkDatabase(ctx)
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		allHealthy = false
	}

	if h.redis != nil {
		redisCheck := h.checkRedis(ctx)
		checks["redis"] = redisCheck
		if redisCheck.Status != "healthy" {
			allHealthy = false
		}
	}

	status := "healthy"
	statusCode := fiber.StatusOK
	if !allHealthy {
		status = "unhealthy"
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(dto.NewHealthResponse(status, checks))
}

func (h *HealthHandler) checkDatabase(ctx context.Context) dto.Check {
	start := time.Now()

	err := h.db.Ping(ctx)
	latency := time.Since(start)

	if err != nil {
		return dto.NewCheck("unhealthy", formatLatency(latency), err.Error())
	}

	return dto.NewCheck("healthy", formatLatency(latency), "")
}

func (h *HealthHandler) checkRedis(ctx context.Context) dto.Check {
	start := time.Now()

	_, err := h.redis.Ping(ctx).Result()
	latency := time.Since(start)

	if err != nil {
		return dto.NewCheck("unhealthy", formatLatency(latency), err.Error())
	}

	return dto.NewCheck("healthy", formatLatency(latency), "")
}

func formatLatency(d time.Duration) string {
	return fmt.Sprintf("%dms", d.Milliseconds())
}
