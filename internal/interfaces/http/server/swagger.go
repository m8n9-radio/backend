package server

import (
	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"

	_ "hub/docs" // Import generated docs
)

// SetupSwagger configures Swagger UI routes.
func SetupSwagger(app *fiber.App) {
	app.Get("/api/docs/*", swagger.WrapHandler)
}
