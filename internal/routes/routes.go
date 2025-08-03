package routes

import (
	"holding-snapshots/internal/controllers"
	"holding-snapshots/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(app *fiber.App) {
	// Middleware global
	app.Use(middleware.CORS())

	// Grupo de rutas API
	api := app.Group("/api")

	// Controladores
	validationController := controllers.NewValidationController()

	// Rutas públicas (sin autenticación)
	api.Get("/health", validationController.HealthCheck)

	// Rutas protegidas (con autenticación API Key)
	protected := api.Group("", middleware.APIKeyAuth())
	protected.Post("/validate", validationController.ValidateHolding)
}