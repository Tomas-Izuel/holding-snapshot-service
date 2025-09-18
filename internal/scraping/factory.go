package scraping

import (
	"fmt"
	"holding-snapshots/internal/models"
)

// ScrapingFactory crea estrategias de scraping según el tipo
type ScrapingFactory struct {
	strategies map[string]ScrapingStrategy
}

func GetStrategy(typeInvestment *models.TypeInvestment) (ScrapingStrategy, error) {
	switch typeInvestment.Name {
	case CedearsStrategyEnum:
		fmt.Println("CedearsStrategyEnum")
		return &StockStrategy{}, nil
	case CryptoStrategyEnum:
		fmt.Println("CryptoStrategyEnum")
		return &StockStrategy{}, nil
	case StockStrategyEnum:
		fmt.Println("StockStrategyEnum")
		return &StockStrategy{}, nil
	default:
		return nil, fmt.Errorf("no se encontró la estrategia para el tipo de inversión: %s", typeInvestment.Name)
	}
}

func GetAssetData(typeInvestment *models.TypeInvestment, code string) (float64, error) {
	strategy, err := GetStrategy(typeInvestment)
	if err != nil {
		return 0, err
	}

	price, err := strategy.FetchPrice(typeInvestment, code)
	if err != nil {
		return 0, err
	}

	return price, nil
}
