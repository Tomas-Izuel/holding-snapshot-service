package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Configurar servicios de cron
	cronService := services.NewCronService()
	if err := cronService.Start(); err != nil {
		log.Fatalf("❌ Error iniciando servicio de cron: %v", err)
	}

	// Configurar rutas
	routes.SetupRoutes(app, cronService)

	// Mensaje de inicio
	log.Printf("🚀 Servidor iniciando en puerto %s", cfg.Port)
	log.Printf("🌍 Entorno: %s", cfg.Environment)
	log.Printf("📅 Próxima ejecución de cron: %s", cronService.GetNextScheduledRun().Format("2006-01-02 15:04:05 UTC"))

	// Configurar graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("🛑 Cerrando servicios...")
		cronService.Stop()
		log.Println("✅ Servicios cerrados correctamente")
		os.Exit(0)
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
