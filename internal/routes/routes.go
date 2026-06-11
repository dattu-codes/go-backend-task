package routes

import (
	"go-backend-task/internal/handler"
	"go-backend-task/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up the routing configuration for the REST API.
// It applies the RequestID and Logger middlewares to all endpoints under the /users namespace.
func RegisterRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	// All routes in this group will run RequestID first, then custom Logger
	api := app.Group("/users", middleware.RequestID(), middleware.Logger())

	api.Post("/", userHandler.Create)
	api.Get("/:id", userHandler.Get)
	api.Put("/:id", userHandler.Update)
	api.Delete("/:id", userHandler.Delete)
	api.Get("/", userHandler.List)
}
