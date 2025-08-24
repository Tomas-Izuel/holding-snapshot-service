package scraping

import "holding-snapshots/internal/models"

// ScrapingStrategy define la interfaz para las estrategias de scraping
type ScrapingStrategy interface {
	// FetchPrice obtiene el precio de un activo usando la estrategia específica
	FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error)

	// BuildURL construye la URL específica para el scraping
	BuildURL(baseURL, code string) string
}

const (
	CedearsStrategyEnum = "cedears"
	CryptoStrategyEnum  = "crypto"
	StockStrategyEnum   = "stock"
)
