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

// CedearsStrategy implementa la estrategia de scraping para CEDEARs
type CedearsStrategy struct{}

// ScrapingResponse representa la respuesta del scraping para CEDEARs
type CedearsScrapingResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Valid  bool    `json:"valid"`
	Market string  `json:"market,omitempty"`
}

// FetchPrice obtiene el precio de un CEDEAR
func (s *CedearsStrategy) FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error) {
	// Verificar cache primero
	cacheKey := fmt.Sprintf("cedear_price:%s", code)
	ctx := context.Background()
	
	if cachedPrice, err := cache.Get(ctx, cacheKey); err == nil {
		var price float64
		if json.Unmarshal([]byte(cachedPrice), &price) == nil {
			return price, nil
		}
	}

	// Construir URL específica para CEDEARs
	fullURL := s.BuildURL(typeInvestment.ScrapingURL, code)
	
	// Realizar request HTTP con headers específicos para CEDEARs
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando request: %w", err)
	}
	
	// Headers específicos para CEDEARs (si es necesario)
	req.Header.Set("User-Agent", "HoldingSnapshots/1.0")
	req.Header.Set("Accept", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error en request HTTP para CEDEARs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code no exitoso para CEDEARs: %d", resp.StatusCode)
	}

	// Leer y parsear respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta de CEDEARs: %w", err)
	}

	var scrapingResp CedearsScrapingResponse
	err = json.Unmarshal(body, &scrapingResp)
	if err != nil {
		return 0, fmt.Errorf("error parseando JSON de CEDEARs: %w", err)
	}

	if !scrapingResp.Valid {
		return 0, fmt.Errorf("CEDEAR %s no válido según el servicio de scraping", code)
	}

	// Guardar en cache por 10 minutos (CEDEARs pueden necesitar más tiempo)
	priceJSON, _ := json.Marshal(scrapingResp.Price)
	cache.Set(ctx, cacheKey, string(priceJSON), 10*time.Minute)

	return scrapingResp.Price, nil
}

// GetSupportedType devuelve el tipo soportado
func (s *CedearsStrategy) GetSupportedType() string {
	return "CEDEARS"
}

// BuildURL construye la URL específica para CEDEARs
func (s *CedearsStrategy) BuildURL(baseURL, code string) string {
	// Para CEDEARs, agregamos parámetros específicos
	return fmt.Sprintf("%s?symbol=%s&market=cedears&currency=ars", baseURL, code)
}