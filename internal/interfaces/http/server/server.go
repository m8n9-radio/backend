package server

import (
	"context"
	"fmt"

	"hub/internal/logger"

	"github.com/gofiber/fiber/v2"
)

// Server represents thserver.
type Server struct {
	app    *fiber.App
	router *Router
	logger *logger.Logger
}

// NewServer creates a new Server.
func NewServer(
	router *Router,
	logger *logger.Logger,
) *Server {
	app := NewFiberApp()
	router.Setup(app)

	return &Server{
		app:    app,
		router: router,
		logger: logger,
	}
}

// Start starts the HTTP server.
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.logger.Info(fmt.Sprintf("Starting HTTP server on %s", addr))
	return s.app.Listen(addr)
}

// Listen starts the HTTP server (alias for Start).
func (s *Server) Listen(port int) error {
	return s.Start(port)
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

// App returns the underlying Fiber app.
func (s *Server) App() *fiber.App {
	return s.app
}
