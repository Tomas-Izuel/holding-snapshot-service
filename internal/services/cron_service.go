package services

import (
	"log"

	"github.com/gofiber/contrib/fibercron"
	"github.com/gofiber/fiber/v2"
)

type CronService struct {
	scrapingService *ScrapingService
}

// NewCronService crea una nueva instancia del servicio de cron
func NewCronService() *CronService {
	return &CronService{
		scrapingService: NewScrapingService(),
	}
}

// SetupCronJobs configura todos los trabajos cron de la aplicación
func (cs *CronService) SetupCronJobs(app *fiber.App, cronSchedule string) {
	log.Printf("📅 Configurando cron job con schedule: %s", cronSchedule)

	// Configurar fibercron
	app.Use(fibercron.New(fibercron.Config{
		TimeZone: "America/Argentina/Buenos_Aires", // Ajustar según la zona horaria deseada
	}))

	// Cron job para scraping semanal
	// Por defecto: "0 1 * * 0" = Domingos a la 1:00 AM
	app.Get("/cron/scraping", fibercron.New(fibercron.Config{
		TimeZone: "America/Argentina/Buenos_Aires",
	}), func(c *fiber.Ctx) error {
		// Ejecutar en una goroutine para no bloquear el request
		go func() {
			log.Println("🔄 Iniciando cron job de scraping...")
			if err := cs.scrapingService.ExecuteWeeklyScraping(); err != nil {
				log.Printf("❌ Error en cron job de scraping: %v", err)
			} else {
				log.Println("✅ Cron job de scraping completado exitosamente")
			}
		}()

		return c.JSON(fiber.Map{
			"message": "Scraping job iniciado",
		})
	})

	// Endpoint manual para ejecutar scraping (útil para testing)
	app.Post("/api/scraping/execute", func(c *fiber.Ctx) error {
		// Verificar API Key
		apiKey := c.Get("Authorization")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API Key requerida",
			})
		}

		go func() {
			log.Println("🔧 Ejecutando scraping manual...")
			if err := cs.scrapingService.ExecuteWeeklyScraping(); err != nil {
				log.Printf("❌ Error en scraping manual: %v", err)
			}
		}()

		return c.JSON(fiber.Map{
			"message": "Scraping manual iniciado",
		})
	})

	log.Println("✅ Cron jobs configurados correctamente")
}