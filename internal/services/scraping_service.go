package services

import (
	"context"
	"encoding/json"
	"fmt"
	"holding-snapshots/internal/models"
	"holding-snapshots/pkg/cache"
	"holding-snapshots/pkg/database"
	"io"
	"log"
	"net/http"
	"time"
)

type ScrapingService struct{}

// NewScrapingService crea una nueva instancia del servicio de scraping
func NewScrapingService() *ScrapingService {
	return &ScrapingService{}
}

// ScrapingResponse representa la respuesta esperada del endpoint de scraping
type ScrapingResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Valid  bool    `json:"valid"`
}

// ExecuteWeeklyScraping ejecuta el scraping semanal de todos los holdings
func (s *ScrapingService) ExecuteWeeklyScraping() error {
	log.Println("🚀 Iniciando scraping semanal...")

	// Obtener todos los holdings con sus grupos y tipos de inversión
	var holdings []models.Holding
	err := database.DB.Preload("Group.Type").Find(&holdings).Error
	if err != nil {
		log.Printf("❌ Error al obtener holdings: %v", err)
		return err
	}

	if len(holdings) == 0 {
		log.Println("ℹ️ No hay holdings para procesar")
		return nil
	}

	log.Printf("📊 Procesando %d holdings...", len(holdings))

	for _, holding := range holdings {
		err = s.ProcessHolding(&holding)
		if err != nil {
			log.Printf("❌ Error procesando holding %s (%s): %v", holding.Name, holding.Code, err)
			continue
		}
		
		// Pequeña pausa entre requests para ser respetuosos con el servidor
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("✅ Scraping semanal completado")
	return nil
}

// ProcessHolding procesa un holding individual
func (s *ScrapingService) ProcessHolding(holding *models.Holding) error {
	// Obtener el precio actual del activo
	price, err := s.FetchAssetPrice(holding.Group.Type.ScrapingURL, holding.Code)
	if err != nil {
		return fmt.Errorf("error obteniendo precio para %s: %w", holding.Code, err)
	}

	// Crear snapshot
	snapshot := models.Snapshot{
		Price:     price,
		HoldingID: holding.ID,
		Quantity:  holding.Quantity,
		CreatedAt: time.Now(),
	}

	err = database.DB.Create(&snapshot).Error
	if err != nil {
		return fmt.Errorf("error creando snapshot: %w", err)
	}

	// Actualizar holding con nuevos cálculos
	holding.CalculateEarnings(price)
	
	err = database.DB.Save(holding).Error
	if err != nil {
		return fmt.Errorf("error actualizando holding: %w", err)
	}

	log.Printf("📈 Holding actualizado: %s (%s) - Precio: %.2f %s", 
		holding.Name, holding.Code, price, holding.Group.Type.Currency)

	return nil
}

// ValidateHolding valida si un holding existe en la URL de scraping
func (s *ScrapingService) ValidateHolding(name, code, groupID string, quantity float64) (*models.Holding, bool, error) {
	// Obtener el grupo y su tipo de inversión
	var group models.Group
	err := database.DB.Preload("Type").First(&group, "id = ?", groupID).Error
	if err != nil {
		return nil, false, fmt.Errorf("grupo no encontrado: %w", err)
	}

	// Verificar si el activo existe en la URL de scraping
	_, err = s.FetchAssetPrice(group.Type.ScrapingURL, code)
	if err != nil {
		log.Printf("⚠️ Holding %s (%s) no válido: %v", name, code, err)
		return nil, false, nil // No es válido pero no es un error del sistema
	}

	// Crear el holding de respuesta (sin guardarlo en DB)
	holding := &models.Holding{
		Name:     name,
		Code:     code,
		GroupID:  groupID,
		Quantity: quantity,
	}

	return holding, true, nil
}

// FetchAssetPrice obtiene el precio de un activo desde la URL de scraping
func (s *ScrapingService) FetchAssetPrice(scrapingURL, code string) (float64, error) {
	// Verificar cache primero (TTL de 5 minutos para evitar múltiples requests)
	cacheKey := fmt.Sprintf("asset_price:%s", code)
	ctx := context.Background()
	
	if cachedPrice, err := cache.Get(ctx, cacheKey); err == nil {
		var price float64
		if json.Unmarshal([]byte(cachedPrice), &price) == nil {
			return price, nil
		}
	}

	// Construir URL completa
	fullURL := fmt.Sprintf("%s?symbol=%s", scrapingURL, code)
	
	// Realizar request HTTP
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return 0, fmt.Errorf("error en request HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code no exitoso: %d", resp.StatusCode)
	}

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Parsear respuesta JSON
	var scrapingResp ScrapingResponse
	err = json.Unmarshal(body, &scrapingResp)
	if err != nil {
		return 0, fmt.Errorf("error parseando JSON: %w", err)
	}

	if !scrapingResp.Valid {
		return 0, fmt.Errorf("activo no válido según el servicio de scraping")
	}

	// Guardar en cache por 5 minutos
	priceJSON, _ := json.Marshal(scrapingResp.Price)
	cache.Set(ctx, cacheKey, string(priceJSON), 5*time.Minute)

	return scrapingResp.Price, nil
}