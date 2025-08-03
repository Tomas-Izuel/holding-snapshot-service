package services

import (
	"holding-snapshots/internal/config"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	scrapingService *ScrapingService
	cronJob         *cron.Cron
}

// NewCronService crea una nueva instancia del servicio de cron
func NewCronService() *CronService {
	// Crear cron con timezone
	location, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		log.Printf("⚠️ Error cargando timezone, usando UTC: %v", err)
		location = time.UTC
	}

	cronJob := cron.New(cron.WithLocation(location))

	return &CronService{
		scrapingService: NewScrapingService(),
		cronJob:         cronJob,
	}
}

// SetupCronJobs configura todos los trabajos cron de la aplicación
func (cs *CronService) SetupCronJobs(app *fiber.App, cronSchedule string) {
	log.Printf("📅 Configurando cron job con schedule: %s", cronSchedule)

	// Agregar trabajo de scraping semanal
	_, err := cs.cronJob.AddFunc(cronSchedule, func() {
		log.Println("🔄 Iniciando cron job de scraping...")
		if err := cs.scrapingService.ExecuteWeeklyScraping(); err != nil {
			log.Printf("❌ Error en cron job de scraping: %v", err)
		} else {
			log.Println("✅ Cron job de scraping completado exitosamente")
		}
	})

	if err != nil {
		log.Printf("❌ Error configurando cron job: %v", err)
		return
	}

	// Iniciar el cron
	cs.cronJob.Start()

	// Endpoint manual para ejecutar scraping (útil para testing)
	app.Post("/api/scraping/execute", func(c *fiber.Ctx) error {
		// Verificar API Key
		apiKey := c.Get("Authorization")
		if apiKey == "" || apiKey != config.AppConfig.APIKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API Key requerida",
			})
		}

		go func() {
			log.Println("🔧 Ejecutando scraping manual...")
			if err := cs.scrapingService.ExecuteWeeklyScraping(); err != nil {
				log.Printf("❌ Error en scraping manual: %v", err)
			} else {
				log.Println("✅ Scraping manual completado")
			}
		}()

		return c.JSON(fiber.Map{
			"message": "Scraping manual iniciado",
		})
	})

	log.Println("✅ Cron jobs configurados correctamente")
}

// Stop detiene el servicio de cron
func (cs *CronService) Stop() {
	if cs.cronJob != nil {
		cs.cronJob.Stop()
		log.Println("🛑 Cron jobs detenidos")
	}
}
