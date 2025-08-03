package scraping

import (
	"context"
	"encoding/json"
	"fmt"
	"holding-snapshots/internal/models"
	"holding-snapshots/pkg/cache"
	"io"
	"net/http"
	"time"
)

// AccionesStrategy implementa la estrategia de scraping para Acciones
type AccionesStrategy struct{}

// AccionesScrapingResponse representa la respuesta del scraping para Acciones
type AccionesScrapingResponse struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Valid     bool    `json:"valid"`
	Exchange  string  `json:"exchange,omitempty"`
	Volume    int64   `json:"volume,omitempty"`
}

// FetchPrice obtiene el precio de una acción
func (s *AccionesStrategy) FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error) {
	// Verificar cache primero
	cacheKey := fmt.Sprintf("accion_price:%s", code)
	ctx := context.Background()
	
	if cachedPrice, err := cache.Get(ctx, cacheKey); err == nil {
		var price float64
		if json.Unmarshal([]byte(cachedPrice), &price) == nil {
			return price, nil
		}
	}

	// Construir URL específica para Acciones
	fullURL := s.BuildURL(typeInvestment.ScrapingURL, code)
	
	// Realizar request HTTP con configuración para acciones
	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando request para acciones: %w", err)
	}
	
	// Headers específicos para acciones
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
	
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error en request HTTP para acciones: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code no exitoso para acciones: %d", resp.StatusCode)
	}

	// Leer y parsear respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta de acciones: %w", err)
	}

	var scrapingResp AccionesScrapingResponse
	err = json.Unmarshal(body, &scrapingResp)
	if err != nil {
		return 0, fmt.Errorf("error parseando JSON de acciones: %w", err)
	}

	if !scrapingResp.Valid {
		return 0, fmt.Errorf("acción %s no válida según el servicio de scraping", code)
	}

	// Guardar en cache por 5 minutos (acciones necesitan actualizaciones más frecuentes)
	priceJSON, _ := json.Marshal(scrapingResp.Price)
	cache.Set(ctx, cacheKey, string(priceJSON), 5*time.Minute)

	return scrapingResp.Price, nil
}

// GetSupportedType devuelve el tipo soportado
func (s *AccionesStrategy) GetSupportedType() string {
	return "ACCIONES"
}

// BuildURL construye la URL específica para Acciones
func (s *AccionesStrategy) BuildURL(baseURL, code string) string {
	// Para acciones, incluimos exchange y tipo de mercado
	return fmt.Sprintf("%s?symbol=%s&type=stock&exchange=NYSE,NASDAQ", baseURL, code)
}