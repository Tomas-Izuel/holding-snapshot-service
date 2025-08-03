package scraping

import (
	"fmt"
	"strings"
)

// ScrapingFactory crea estrategias de scraping según el tipo
type ScrapingFactory struct {
	strategies map[string]ScrapingStrategy
}

// NewScrapingFactory crea una nueva instancia de la factory
func NewScrapingFactory() *ScrapingFactory {
	factory := &ScrapingFactory{
		strategies: make(map[string]ScrapingStrategy),
	}

	// Registrar todas las estrategias disponibles
	factory.RegisterStrategy(&CedearsStrategy{})
	factory.RegisterStrategy(&AccionesStrategy{})
	factory.RegisterStrategy(&CryptoStrategy{})

	return factory
}

// RegisterStrategy registra una nueva estrategia en la factory
func (f *ScrapingFactory) RegisterStrategy(strategy ScrapingStrategy) {
	supportedType := strings.ToUpper(strategy.GetSupportedType())
	f.strategies[supportedType] = strategy
}

// GetStrategy devuelve la estrategia apropiada según el nombre del grupo
func (f *ScrapingFactory) GetStrategy(groupName string) (ScrapingStrategy, error) {
	// Normalizar el nombre del grupo a mayúsculas para la comparación
	normalizedName := strings.ToUpper(strings.TrimSpace(groupName))

	// Mapeo de nombres de grupos a tipos de estrategia
	var strategyType string
	switch {
	case strings.Contains(normalizedName, "CEDEAR"):
		strategyType = "CEDEARS"
	case strings.Contains(normalizedName, "ACCION"):
		strategyType = "ACCIONES"
	case strings.Contains(normalizedName, "CRYPTO") || strings.Contains(normalizedName, "CRIPTO"):
		strategyType = "CRYPTO"
	default:
		// Intentar mapeo directo si no coincide con los patrones
		strategyType = normalizedName
	}

	strategy, exists := f.strategies[strategyType]
	if !exists {
		return nil, fmt.Errorf("no existe estrategia de scraping para el tipo de grupo: %s", groupName)
	}

	return strategy, nil
}

// GetAvailableStrategies devuelve la lista de estrategias disponibles
func (f *ScrapingFactory) GetAvailableStrategies() []string {
	strategies := make([]string, 0, len(f.strategies))
	for strategyType := range f.strategies {
		strategies = append(strategies, strategyType)
	}
	return strategies
}

// HasStrategy verifica si existe una estrategia para el tipo dado
func (f *ScrapingFactory) HasStrategy(groupName string) bool {
	_, err := f.GetStrategy(groupName)
	return err == nil
}

// GetStrategyInfo devuelve información sobre la estrategia para un grupo
func (f *ScrapingFactory) GetStrategyInfo(groupName string) (string, error) {
	strategy, err := f.GetStrategy(groupName)
	if err != nil {
		return "", err
	}
	return strategy.GetSupportedType(), nil
}

// ValidateGroupName verifica si un nombre de grupo es válido y sugiere el tipo
func (f *ScrapingFactory) ValidateGroupName(groupName string) (bool, string, string) {
	strategyType, err := f.GetStrategyInfo(groupName)
	if err != nil {
		// Sugerir tipos disponibles
		suggestions := strings.Join(f.GetAvailableStrategies(), ", ")
		return false, "", fmt.Sprintf("Tipos disponibles: %s", suggestions)
	}
	return true, strategyType, fmt.Sprintf("Grupo '%s' será procesado con estrategia %s", groupName, strategyType)
}
