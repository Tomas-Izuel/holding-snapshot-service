package scraping

import "holding-snapshots/internal/models"

// ScrapingStrategy define la interfaz para las estrategias de scraping
type ScrapingStrategy interface {
	// FetchPrice obtiene el precio de un activo usando la estrategia específica
	FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error)

	// GetSupportedType devuelve el tipo de inversión que soporta esta estrategia
	GetSupportedType() string

	// BuildURL construye la URL específica para el scraping
	BuildURL(baseURL, code string) string
}
