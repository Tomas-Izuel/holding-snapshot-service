package services

import (
	"fmt"
	"log"
	"time"

	"holding-snapshots/internal/models"
	"holding-snapshots/internal/scraping"
	"holding-snapshots/pkg/database"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	cron            *cron.Cron
	scrapingService *ScrapingService
}

// NewCronService crea una nueva instancia del servicio de cron
func NewCronService() *CronService {
	// Crear cron con timezone UTC
	c := cron.New(cron.WithLocation(time.UTC))

	return &CronService{
		cron:            c,
		scrapingService: NewScrapingService(),
	}
}

// Start inicia el servicio de cron
func (cs *CronService) Start() error {
	log.Println("🚀 Iniciando servicio de cron...")

	// Programar cronjob para ejecutarse los domingos a las 3:00 AM UTC
	// Cron expression: "0 3 * * 0" (minuto 0, hora 3, cualquier día del mes, cualquier mes, domingo)
	_, err := cs.cron.AddFunc("0 3 * * 0", cs.ExecuteWeeklyScraping)
	if err != nil {
		return fmt.Errorf("error programando cronjob semanal: %w", err)
	}

	log.Println("📅 Cronjob programado para ejecutarse los domingos a las 3:00 AM UTC")

	// Iniciar el cron
	cs.cron.Start()
	log.Println("✅ Servicio de cron iniciado correctamente")

	return nil
}

// Stop detiene el servicio de cron
func (cs *CronService) Stop() {
	log.Println("⏹️ Deteniendo servicio de cron...")
	cs.cron.Stop()
	log.Println("✅ Servicio de cron detenido")
}

// ExecuteWeeklyScraping ejecuta el scraping semanal de todos los assets
func (cs *CronService) ExecuteWeeklyScraping() {
	log.Println("🚀 Iniciando scraping semanal de assets...")
	startTime := time.Now()

	// Obtener todos los assets válidos con su tipo de inversión
	assets, err := cs.getAllValidAssets()
	if err != nil {
		log.Printf("❌ Error obteniendo assets: %v", err)
		return
	}

	if len(assets) == 0 {
		log.Println("ℹ️ No hay assets válidos para procesar")
		return
	}

	log.Printf("📊 Procesando %d assets...", len(assets))

	successCount := 0
	errorCount := 0

	// Procesar cada asset
	for _, asset := range assets {
		err := cs.processAsset(&asset)
		if err != nil {
			log.Printf("❌ Error procesando asset %s (%s): %v", asset.Name, asset.Code, err)
			errorCount++
		} else {
			successCount++
			log.Printf("✅ Asset procesado exitosamente: %s (%s) - Precio: %.2f",
				asset.Name, asset.Code, asset.LastPrice)
		}

		// Pausa pequeña entre requests para ser respetuosos
		time.Sleep(200 * time.Millisecond)
	}

	duration := time.Since(startTime)
	log.Printf("🏁 Scraping semanal completado en %v - Éxitos: %d, Errores: %d",
		duration, successCount, errorCount)
}

// getAllValidAssets obtiene todos los assets válidos con su tipo de inversión
func (cs *CronService) getAllValidAssets() ([]models.Asset, error) {
	var assets []models.Asset

	// Usar el nombre del campo del struct para evitar problemas de naming
	err := database.DB.Preload("Type").Where(&models.Asset{IsValid: true}).Find(&assets).Error
	if err != nil {
		return nil, fmt.Errorf("error obteniendo assets válidos: %w", err)
	}

	return assets, nil
}

// processAsset procesa un asset individual: scrapea precio y crea snapshots
func (cs *CronService) processAsset(asset *models.Asset) error {
	log.Printf("🔍 Procesando asset: %s (%s)", asset.Name, asset.Code)

	// Scrapear el precio actual del asset
	price, err := cs.scrapeAssetPrice(asset)
	if err != nil {
		return fmt.Errorf("error scrapeando precio: %w", err)
	}

	// Actualizar el lastPrice del asset para optimización futura
	err = cs.updateAssetLastPrice(asset, price)
	if err != nil {
		return fmt.Errorf("error actualizando precio del asset: %w", err)
	}

	// Crear snapshots para todos los holdings de este asset
	err = cs.createSnapshotsForAsset(asset, price)
	if err != nil {
		return fmt.Errorf("error creando snapshots: %w", err)
	}

	return nil
}

// scrapeAssetPrice scrapea el precio actual de un asset usando la estrategia correcta
func (cs *CronService) scrapeAssetPrice(asset *models.Asset) (float64, error) {
	// Obtener la estrategia correcta según el tipo de inversión
	strategy, err := scraping.GetStrategy(&asset.Type)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo estrategia para tipo %s: %w", asset.Type.Name, err)
	}

	// Scrapear el precio
	price, err := strategy.FetchPrice(&asset.Type, asset.Code)
	if err != nil {
		return 0, fmt.Errorf("error fetching price con estrategia %s: %w", asset.Type.Name, err)
	}

	log.Printf("💰 Precio scrapeado para %s (%s): %.2f %s",
		asset.Name, asset.Code, price, asset.Type.Currency)

	return price, nil
}

// updateAssetLastPrice actualiza el lastPrice del asset en la base de datos
func (cs *CronService) updateAssetLastPrice(asset *models.Asset, price float64) error {
	asset.LastPrice = price

	err := database.DB.Save(asset).Error
	if err != nil {
		return fmt.Errorf("error guardando asset actualizado: %w", err)
	}

	return nil
}

// createSnapshotsForAsset crea snapshots para todos los holdings de un asset
func (cs *CronService) createSnapshotsForAsset(asset *models.Asset, currentPrice float64) error {
	// Obtener todos los holdings de este asset
	var holdings []models.Holding
	err := database.DB.Where("asset_id = ?", asset.ID).Find(&holdings).Error
	if err != nil {
		return fmt.Errorf("error obteniendo holdings para asset %s: %w", asset.ID, err)
	}

	if len(holdings) == 0 {
		log.Printf("ℹ️ No hay holdings para el asset %s (%s)", asset.Name, asset.Code)
		return nil
	}

	log.Printf("📸 Creando %d snapshots para asset %s (%s)", len(holdings), asset.Name, asset.Code)

	// Crear snapshot para cada holding
	for _, holding := range holdings {
		snapshot := models.Snapshot{
			Price:     currentPrice,
			HoldingID: holding.ID,
			Quantity:  holding.Quantity,
			CreatedAt: time.Now(),
		}

		err := database.DB.Create(&snapshot).Error
		if err != nil {
			log.Printf("⚠️ Error creando snapshot para holding %s: %v", holding.ID, err)
			continue
		}

		// Actualizar earnings del holding
		err = cs.updateHoldingEarnings(&holding, currentPrice)
		if err != nil {
			log.Printf("⚠️ Error actualizando earnings para holding %s: %v", holding.ID, err)
			continue
		}
	}

	return nil
}

// updateHoldingEarnings actualiza las ganancias del holding basado en el precio actual
func (cs *CronService) updateHoldingEarnings(holding *models.Holding, currentPrice float64) error {
	// Obtener el snapshot más reciente anterior a este para calcular earnings
	var previousSnapshot models.Snapshot
	err := database.DB.Where("holding_id = ?", holding.ID).
		Order("created_at DESC").
		Offset(1). // Saltar el snapshot que acabamos de crear
		First(&previousSnapshot).Error

	if err != nil {
		// Si no hay snapshot anterior, no podemos calcular earnings
		log.Printf("ℹ️ No hay snapshot anterior para holding %s, earnings se mantienen en 0", holding.ID)
		return nil
	}

	// Calcular earnings usando el método del modelo
	holding.CalculateEarnings(currentPrice, previousSnapshot.Price)

	// Guardar el holding actualizado
	err = database.DB.Save(holding).Error
	if err != nil {
		return fmt.Errorf("error guardando holding actualizado: %w", err)
	}

	log.Printf("📈 Earnings actualizados para holding %s: %.2f (%.2f%%)",
		holding.ID, holding.Earnings, holding.RelativeEarnings)

	return nil
}

// ExecuteManualScraping permite ejecutar el scraping manualmente para testing
func (cs *CronService) ExecuteManualScraping() {
	log.Println("🔧 Ejecutando scraping manual...")
	cs.ExecuteWeeklyScraping()
}

// GetNextScheduledRun obtiene la próxima ejecución programada
func (cs *CronService) GetNextScheduledRun() time.Time {
	entries := cs.cron.Entries()
	if len(entries) > 0 {
		return entries[0].Next
	}
	return time.Time{}
}

// GetCronStatus obtiene el estado del servicio de cron
func (cs *CronService) GetCronStatus() map[string]interface{} {
	entries := cs.cron.Entries()

	status := map[string]interface{}{
		"running":        len(entries) > 0,
		"total_jobs":     len(entries),
		"next_execution": nil,
	}

	if len(entries) > 0 {
		status["next_execution"] = entries[0].Next.Format("2006-01-02 15:04:05 UTC")
	}

	return status
}
