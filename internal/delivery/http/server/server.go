package server

import (
	"context"
	"fmt"

	"hub/internal/delivery/http/handler"
	"hub/internal/delivery/http/middleware"
	"hub/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Server interface {
		Listen(port int) error
		Shutdown(ctx context.Context) error
	}
	server struct {
		app    *fiber.App
		logger *logger.Logger
		pool   *pgxpool.Pool

		trackHandler      handler.TrackHandler
		reactionHandler   handler.ReactionHandler
		radioHandler      handler.RadioHandler
		statisticsHandler handler.StatisticsHandler
	}
)

func NewServer(logger *logger.Logger, pool *pgxpool.Pool, trackHandler handler.TrackHandler, reactionHandler handler.ReactionHandler, radioHandler handler.RadioHandler, statisticsHandler handler.StatisticsHandler) Server {
	app := FiberApplication()

	server := &server{
		app:    app,
		logger: logger,
		pool:   pool,

		trackHandler:      trackHandler,
		reactionHandler:   reactionHandler,
		radioHandler:      radioHandler,
		statisticsHandler: statisticsHandler,
	}

	server.useMiddlewares()
	server.useRouters()

	return server
}

func (s *server) useMiddlewares() {
	s.app.Use(middleware.LoggingMiddleware(s.logger))
}

func (s *server) useRouters() {
	s.useTrackRoute()
	s.useReactionRoute()
	s.useRadioRoute()
	s.useStatisticsRoute()
}

func (s *server) Listen(port int) error {
	return s.app.Listen(
		fmt.Sprintf("0.0.0.0:%d", port),
	)
}

func (s *server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")
	return s.app.ShutdownWithContext(ctx)
}
