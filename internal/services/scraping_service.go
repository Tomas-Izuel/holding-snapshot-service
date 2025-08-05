package services

import (
	"context"
	"encoding/json"
	"fmt"
	"holding-snapshots/internal/models"
	"holding-snapshots/internal/scraping"
	"holding-snapshots/pkg/cache"
	"holding-snapshots/pkg/database"
	"log"
	"time"
)

type ScrapingService struct {
	factory *scraping.ScrapingFactory
}

// NewScrapingService crea una nueva instancia del servicio de scraping
func NewScrapingService() *ScrapingService {
	return &ScrapingService{
		factory: scraping.NewScrapingFactory(),
	}
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
	// Obtener el precio actual del activo usando la factory
	price, err := s.FetchAssetPrice(&holding.Group.Type, holding.Code, holding.Group.Name)
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

// ValidatedHoldingCache representa la estructura de datos para el cache de holdings validados
type ValidatedHoldingCache struct {
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	TypeID   string    `json:"type_id"`
	TypeName string    `json:"type_name"`
	Valid    bool      `json:"valid"`
	CachedAt time.Time `json:"cached_at"`
}

// ValidateHoldingWithType valida si un holding existe usando directamente el tipo de inversión
// Este método es útil para validar holdings antes de crear un grupo
func (s *ScrapingService) ValidateHoldingWithType(name, code string, typeInvestment *models.TypeInvestment, groupName string, quantity float64) (*models.Holding, bool, error) {
	// Crear clave de cache basada en tipo de inversión y código
	cacheKey := fmt.Sprintf("validated_holding:%s:%s", typeInvestment.ID, code)
	ctx := context.Background()

	// Verificar cache primero
	if cachedData, err := cache.Get(ctx, cacheKey); err == nil {
		var cachedHolding ValidatedHoldingCache
		if json.Unmarshal([]byte(cachedData), &cachedHolding) == nil {
			log.Printf("📦 Holding %s (%s) encontrado en cache para tipo %s - Válido: %v",
				name, code, cachedHolding.TypeName, cachedHolding.Valid)

			if cachedHolding.Valid {
				// Crear el holding usando los datos cacheados
				holding := &models.Holding{
					Name:     name,
					Code:     code,
					Quantity: quantity,
				}
				return holding, true, nil
			} else {
				// Si el cache dice que no es válido, devolver false sin hacer scraping
				return nil, false, nil
			}
		}
	}

	log.Printf("🔍 Holding %s (%s) no encontrado en cache, realizando validación por scraping...", name, code)

	// Debug: Log antes de llamar fetchAssetPriceWithType
	log.Printf("🚀 ANTES de fetchAssetPriceWithType - TypeInvestment.ScrapingURL: '%s'", typeInvestment.ScrapingURL)

	// No está en cache, verificar si el activo existe usando la estrategia apropiada
	_, err := s.fetchAssetPriceWithType(typeInvestment, code, groupName)

	// Crear estructura para cache
	cacheData := ValidatedHoldingCache{
		Name:     name,
		Code:     code,
		TypeID:   typeInvestment.ID,
		TypeName: typeInvestment.Name,
		Valid:    err == nil,
		CachedAt: time.Now(),
	}

	// Guardar resultado en cache (tanto si es válido como si no lo es)
	// TTL de 24 horas para holdings válidos, 2 horas para inválidos
	cacheTTL := 2 * time.Hour
	if cacheData.Valid {
		cacheTTL = 24 * time.Hour
	}

	if cacheJSON, jsonErr := json.Marshal(cacheData); jsonErr == nil {
		cache.Set(ctx, cacheKey, string(cacheJSON), cacheTTL)
		log.Printf("💾 Resultado de validación guardado en cache para %s (%s) - Válido: %v, TTL: %v",
			name, code, cacheData.Valid, cacheTTL)
	}

	if err != nil {
		log.Printf("⚠️ Holding %s (%s) no válido: %v", name, code, err)
		return nil, false, nil // No es válido pero no es un error del sistema
	}

	// Crear el holding de respuesta (sin guardarlo en DB)
	holding := &models.Holding{
		Name:     name,
		Code:     code,
		Quantity: quantity,
	}

	return holding, true, nil
}

// ValidateHolding valida si un holding existe en la URL de scraping con cache Redis
// Este método requiere que el grupo ya exista para obtener el tipo de inversión
func (s *ScrapingService) ValidateHolding(name, code, groupID string, quantity float64) (*models.Holding, bool, error) {
	// Obtener el grupo y su tipo de inversión
	var group models.Group
	err := database.DB.Preload("Type").First(&group, "id = ?", groupID).Error
	if err != nil {
		return nil, false, fmt.Errorf("grupo no encontrado: %w", err)
	}

	// Usar el nuevo método que no depende de grupo existente
	holding, isValid, err := s.ValidateHoldingWithType(name, code, &group.Type, group.Name, quantity)
	if err != nil {
		return nil, false, err
	}

	// Si es válido, agregar el GroupID al holding
	if isValid && holding != nil {
		holding.GroupID = groupID
	}

	return holding, isValid, nil
}

// ClearValidatedHoldingCache elimina el cache de un holding específico
func (s *ScrapingService) ClearValidatedHoldingCache(typeID, code string) error {
	cacheKey := fmt.Sprintf("validated_holding:%s:%s", typeID, code)
	ctx := context.Background()

	err := cache.Delete(ctx, cacheKey)
	if err != nil {
		return fmt.Errorf("error eliminando cache para holding %s del tipo %s: %w", code, typeID, err)
	}

	log.Printf("🗑️ Cache eliminado para holding %s del tipo %s", code, typeID)
	return nil
}

// GetValidatedHoldingFromCache obtiene información de un holding desde el cache
func (s *ScrapingService) GetValidatedHoldingFromCache(typeID, code string) (*ValidatedHoldingCache, bool, error) {
	cacheKey := fmt.Sprintf("validated_holding:%s:%s", typeID, code)
	ctx := context.Background()

	cachedData, err := cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, false, nil // No está en cache, no es un error
	}

	var cachedHolding ValidatedHoldingCache
	if err := json.Unmarshal([]byte(cachedData), &cachedHolding); err != nil {
		return nil, false, fmt.Errorf("error deserializando cache: %w", err)
	}

	return &cachedHolding, true, nil
}

// GetValidationCacheStats obtiene estadísticas del cache de validaciones
func (s *ScrapingService) GetValidationCacheStats() map[string]interface{} {
	stats := map[string]interface{}{
		"cache_pattern": "validated_holding:*",
		"description":   "Cache de holdings validados por tipo de inversión",
		"ttl_valid":     "24 horas para holdings válidos",
		"ttl_invalid":   "2 horas para holdings inválidos",
	}
	return stats
}

// FetchAssetPrice obtiene el precio de un activo usando la estrategia apropiada
// Este método requiere un grupo existente para obtener el tipo de inversión
func (s *ScrapingService) FetchAssetPrice(typeInvestment *models.TypeInvestment, code, groupName string) (float64, error) {
	return s.fetchAssetPriceWithType(typeInvestment, code, groupName)
}

// fetchAssetPriceWithType obtiene el precio usando directamente el tipo de inversión
// Este método no depende de la existencia de un grupo
func (s *ScrapingService) fetchAssetPriceWithType(typeInvestment *models.TypeInvestment, code, groupName string) (float64, error) {
	log.Printf("🔍 fetchAssetPriceWithType llamado - TypeInvestment: %+v, Code: %s, GroupName: %s", typeInvestment, code, groupName)

	// Obtener la estrategia apropiada según el nombre del grupo
	strategy, err := s.factory.GetStrategy(groupName)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo estrategia de scraping para grupo %s: %w", groupName, err)
	}

	log.Printf("🔍 Estrategia obtenida: %s", strategy.GetSupportedType())

	// Usar la estrategia para obtener el precio
	price, err := strategy.FetchPrice(typeInvestment, code)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo precio usando estrategia %s: %w", strategy.GetSupportedType(), err)
	}

	log.Printf("💰 Precio obtenido usando estrategia %s para %s: %.2f %s",
		strategy.GetSupportedType(), code, price, typeInvestment.Currency)

	return price, nil
}
