package server

import (
	"hub/internal/interfaces/http/handler"
	"hub/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// Router configures all HTTP routes.
type Router struct {
	trackHandler      *handler.TrackHandler
	reactionHandler   *handler.ReactionHandler
	radioHandler      *handler.RadioHandler
	statisticsHandler *handler.StatisticsHandler
	healthHandler     *handler.HealthHandler
}

// NewRouter creates a new Router.
func NewRouter(
	trackHandler *handler.TrackHandler,
	reactionHandler *handler.ReactionHandler,
	radioHandler *handler.RadioHandler,
	statisticsHandler *handler.StatisticsHandler,
	healthHandler *handler.HealthHandler,
) *Router {
	return &Router{
		trackHandler:      trackHandler,
		reactionHandler:   reactionHandler,
		radioHandler:      radioHandler,
		statisticsHandler: statisticsHandler,
		healthHandler:     healthHandler,
	}
}

// Setup configures all routes on the Fiber app.
func (r *Router) Setup(app *fiber.App) {
	// Health check
	app.Get("/health", r.healthHandler.Health)

	// Track routes
	app.Get("/tracks/:id", r.trackHandler.Get)
	app.Post("/tracks", middleware.ValidateTrackRequest(), r.trackHandler.Upsert)

	// Reaction routes
	app.Post("/tracks/:trackId/like", r.reactionHandler.Like)
	app.Post("/tracks/:trackId/dislike", r.reactionHandler.Dislike)
	app.Get("/tracks/:trackId/reaction", r.reactionHandler.Check)

	// Radio routes
	app.Get("/radio/info", r.radioHandler.GetInfo)
	app.Get("/radio/listeners", r.radioHandler.GetListen)

	// Statistics routes
	app.Get("/radio/statistics", r.statisticsHandler.GetStatistics)
}
