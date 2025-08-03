package scraping

import (
	"context"
	"encoding/json"
	"fmt"
	"holding-snapshots/internal/models"
	"holding-snapshots/pkg/cache"
	"io"
	"net/http"
	"strings"
	"time"
)

// CryptoStrategy implementa la estrategia de scraping para Criptomonedas
type CryptoStrategy struct{}

// CryptoScrapingResponse representa la respuesta del scraping para Crypto
type CryptoScrapingResponse struct {
	Symbol      string  `json:"symbol"`
	Price       float64 `json:"price"`
	Valid       bool    `json:"valid"`
	PriceUSD    float64 `json:"price_usd,omitempty"`
	MarketCap   int64   `json:"market_cap,omitempty"`
	Volume24h   float64 `json:"volume_24h,omitempty"`
}

// FetchPrice obtiene el precio de una criptomoneda
func (s *CryptoStrategy) FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error) {
	// Verificar cache primero
	cacheKey := fmt.Sprintf("crypto_price:%s", strings.ToUpper(code))
	ctx := context.Background()
	
	if cachedPrice, err := cache.Get(ctx, cacheKey); err == nil {
		var price float64
		if json.Unmarshal([]byte(cachedPrice), &price) == nil {
			return price, nil
		}
	}

	// Construir URL específica para Crypto
	fullURL := s.BuildURL(typeInvestment.ScrapingURL, code)
	
	// Realizar request HTTP con configuración para crypto
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando request para crypto: %w", err)
	}
	
	// Headers específicos para APIs de crypto
	req.Header.Set("User-Agent", "HoldingSnapshots-CryptoBot/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error en request HTTP para crypto: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code no exitoso para crypto: %d", resp.StatusCode)
	}

	// Leer y parsear respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta de crypto: %w", err)
	}

	var scrapingResp CryptoScrapingResponse
	err = json.Unmarshal(body, &scrapingResp)
	if err != nil {
		return 0, fmt.Errorf("error parseando JSON de crypto: %w", err)
	}

	if !scrapingResp.Valid {
		return 0, fmt.Errorf("criptomoneda %s no válida según el servicio de scraping", code)
	}

	// Para crypto, usar el precio en USD si está disponible
	price := scrapingResp.Price
	if scrapingResp.PriceUSD > 0 && typeInvestment.Currency == "USD" {
		price = scrapingResp.PriceUSD
	}

	// Guardar en cache por 3 minutos (crypto cambia muy rápido)
	priceJSON, _ := json.Marshal(price)
	cache.Set(ctx, cacheKey, string(priceJSON), 3*time.Minute)

	return price, nil
}

// GetSupportedType devuelve el tipo soportado
func (s *CryptoStrategy) GetSupportedType() string {
	return "CRYPTO"
}

// BuildURL construye la URL específica para Crypto
func (s *CryptoStrategy) BuildURL(baseURL, code string) string {
	// Para crypto, el símbolo debe estar en mayúsculas y agregar parámetros de conversión
	upperCode := strings.ToUpper(code)
	return fmt.Sprintf("%s?symbol=%s&convert=USD,ARS&include_market_data=true", baseURL, upperCode)
}