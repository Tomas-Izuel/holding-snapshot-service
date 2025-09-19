package controllers

import (
	"holding-snapshots/internal/services"
	"holding-snapshots/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type CronController struct {
	cronService *services.CronService
}

// NewCronController crea una nueva instancia del controlador de cron
func NewCronController(cronService *services.CronService) *CronController {
	return &CronController{
		cronService: cronService,
	}
}

// GetCronStatus obtiene el estado actual del servicio de cron
func (cc *CronController) GetCronStatus(c *fiber.Ctx) error {
	status := cc.cronService.GetCronStatus()

	return utils.SuccessResponse(c, "Estado del cron obtenido exitosamente", status)
}

// ExecuteManualScraping ejecuta el scraping de forma manual
func (cc *CronController) ExecuteManualScraping(c *fiber.Ctx) error {
	// Ejecutar en background para no bloquear la respuesta HTTP
	go cc.cronService.ExecuteManualScraping()

	return utils.SuccessResponse(c, "Scraping manual iniciado en segundo plano", fiber.Map{
		"message": "El scraping se está ejecutando. Revisa los logs del servidor para ver el progreso.",
	})
}

// GetNextScheduledRun obtiene la próxima ejecución programada
func (cc *CronController) GetNextScheduledRun(c *fiber.Ctx) error {
	nextRun := cc.cronService.GetNextScheduledRun()

	data := fiber.Map{
		"next_execution":  nextRun.Format("2006-01-02 15:04:05 UTC"),
		"timezone":        "UTC",
		"cron_expression": "0 3 * * 0", // Domingos a las 3:00 AM
		"description":     "Se ejecuta todos los domingos a las 3:00 AM UTC",
	}

	return utils.SuccessResponse(c, "Próxima ejecución obtenida exitosamente", data)
}

// GetCronInfo obtiene información general sobre el funcionamiento del cron
func (cc *CronController) GetCronInfo(c *fiber.Ctx) error {
	status := cc.cronService.GetCronStatus()
	nextRun := cc.cronService.GetNextScheduledRun()

	info := fiber.Map{
		"service_name":    "Weekly Asset Scraping Service",
		"description":     "Servicio que scrapea precios de assets y crea snapshots de holdings cada domingo",
		"schedule":        "Domingos a las 3:00 AM UTC",
		"cron_expression": "0 3 * * 0",
		"status":          status,
		"next_execution":  nextRun.Format("2006-01-02 15:04:05 UTC"),
		"features": []string{
			"Scraping automático de precios de assets",
			"Actualización del campo lastPrice para optimización",
			"Creación de snapshots para todos los holdings",
			"Cálculo automático de earnings",
			"Logging detallado del proceso",
			"Ejecución manual via API",
		},
		"endpoints": []fiber.Map{
			{
				"method":      "GET",
				"path":        "/api/admin/cron/status",
				"description": "Obtener estado del servicio de cron",
			},
			{
				"method":      "GET",
				"path":        "/api/admin/cron/next",
				"description": "Obtener próxima ejecución programada",
			},
			{
				"method":      "POST",
				"path":        "/api/admin/cron/execute",
				"description": "Ejecutar scraping manual",
			},
			{
				"method":      "GET",
				"path":        "/api/admin/cron/info",
				"description": "Obtener información general del servicio",
			},
		},
	}

	return utils.SuccessResponse(c, "Información del servicio de cron", info)
}
