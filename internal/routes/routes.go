package routes

import (
	"holding-snapshots/internal/controllers"
	"holding-snapshots/internal/middleware"
	"holding-snapshots/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(app *fiber.App, cronService *services.CronService) {
	// Middleware global
	app.Use(middleware.CORS())

	// Grupo de rutas API
	api := app.Group("/api")

	// Controladores
	validationController := controllers.NewValidationController()
	cronController := controllers.NewCronController(cronService)

	// Rutas públicas (sin autenticación)
	api.Get("/health", validationController.HealthCheck)

	// Rutas protegidas (con autenticación API Key)
	protected := api.Group("", middleware.APIKeyAuth())
	protected.Post("/validate", validationController.ValidateHolding)

	// Rutas de administración del cron
	admin := protected.Group("/admin")
	setupCronRoutes(admin, cronController)
}

// setupCronRoutes configura las rutas relacionadas con el servicio de cron
func setupCronRoutes(router fiber.Router, cronController *controllers.CronController) {
	// Obtener estado del cron
	router.Get("/cron/status", cronController.GetCronStatus)

	// Ejecutar scraping manual
	router.Post("/cron/execute", cronController.ExecuteManualScraping)

	// Obtener próxima ejecución programada
	router.Get("/cron/next", cronController.GetNextScheduledRun)

	// Obtener información general del servicio de cron
	router.Get("/cron/info", cronController.GetCronInfo)
}
