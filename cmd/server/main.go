package main

import (
	"log"

	"holding-snapshots/internal/config"
	"holding-snapshots/internal/routes"
	"holding-snapshots/internal/services"
	"holding-snapshots/pkg/cache"
	"holding-snapshots/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a la base de datos
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}

	// Habilitar extensión UUID en PostgreSQL
	if err := database.EnableUUIDExtension(); err != nil {
		log.Printf("⚠️ Advertencia: Error habilitando extensión UUID: %v", err)
	}

	// Ejecutar migraciones
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("❌ Error ejecutando migraciones: %v", err)
	}

	// Conectar a Redis
	if err := cache.Connect(cfg.RedisURL); err != nil {
		log.Fatalf("❌ Error conectando a Redis: %v", err)
	}

	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName:      "Holding Snapshots Service",
		ServerHeader: "Holding-Snapshots",
		ErrorHandler: errorHandler,
	})

	// Middleware global
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))
	app.Use(recover.New())

	// Configurar rutas
	routes.SetupRoutes(app)

	// Configurar servicios de cron
	cronService := services.NewCronService()
	cronService.SetupCronJobs(app, cfg.ScrapingCronSchedule)

	// Mensaje de inicio
	log.Printf("🚀 Servidor iniciando en puerto %s", cfg.Port)
	log.Printf("🌍 Entorno: %s", cfg.Environment)
	log.Printf("📅 Cron schedule: %s", cfg.ScrapingCronSchedule)

	// Configurar graceful shutdown
	defer func() {
		log.Println("🛑 Cerrando servicios...")
		cronService.Stop()
	}()

	// Iniciar servidor
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

// errorHandler maneja errores globales de la aplicación
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Printf("❌ Error: %v", err)

	return c.Status(code).JSON(fiber.Map{
		"error":   "Error interno del servidor",
		"message": err.Error(),
	})
}
