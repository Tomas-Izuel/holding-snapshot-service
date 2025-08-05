package scraping

import (
	"context"
	"encoding/json"
	"fmt"
	"holding-snapshots/internal/models"
	"holding-snapshots/pkg/cache"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AccionesStrategy implementa la estrategia de scraping para Acciones desde Yahoo Finance
type AccionesStrategy struct{}

// YahooFinanceData representa los datos estructurados de Yahoo Finance
type YahooFinanceData struct {
	RegularMarketPrice struct {
		Raw float64 `json:"raw"`
		Fmt string  `json:"fmt"`
	} `json:"regularMarketPrice"`
}

// FetchPrice obtiene el precio de una acción desde Yahoo Finance HTML
func (s *AccionesStrategy) FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error) {
	fmt.Print(typeInvestment)
	fmt.Printf("🎯 AccionesStrategy.FetchPrice llamado con code=%s, scrapingURL=%s\n", code, typeInvestment.ScrapingURL)
	// Verificar cache primero
	cacheKey := fmt.Sprintf("accion_price:%s", code)
	ctx := context.Background()

	if cachedPrice, err := cache.Get(ctx, cacheKey); err == nil {
		var price float64
		if json.Unmarshal([]byte(cachedPrice), &price) == nil {
			return price, nil
		}
	}

	// Construir URL usando la configuración del tipo de inversión
	targetURL := s.BuildURL(typeInvestment.ScrapingURL, code)

	// Debug: Log de la URL construida
	fmt.Printf("🔗 URL construida para %s: %s (base: %s)\n", code, targetURL, typeInvestment.ScrapingURL)

	// Realizar request HTTP con headers que simulan un navegador
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando request para acciones: %w", err)
	}

	// Headers para simular un navegador real
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error en request HTTP para acciones: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code no exitoso para acciones: %d", resp.StatusCode)
	}

	// Leer el HTML
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta de acciones: %w", err)
	}

	htmlContent := string(body)

	// Extraer el precio usando múltiples estrategias
	price, err := s.extractPriceFromHTML(htmlContent, code)
	if err != nil {
		return 0, fmt.Errorf("error extrayendo precio para %s: %w", code, err)
	}

	// Guardar en cache por 5 minutos
	priceJSON, _ := json.Marshal(price)
	cache.Set(ctx, cacheKey, string(priceJSON), 5*time.Minute)

	return price, nil
}

// extractPriceFromHTML extrae el precio del HTML de Yahoo Finance usando múltiples estrategias
func (s *AccionesStrategy) extractPriceFromHTML(htmlContent, symbol string) (float64, error) {
	// Estrategia 1: Buscar en el elemento fin-streamer con data-symbol específico
	price, err := s.extractFromFinStreamer(htmlContent, symbol)
	if err == nil {
		return price, nil
	}

	// Estrategia 2: Buscar en objetos JSON embebidos con regularMarketPrice
	price, err = s.extractFromJSON(htmlContent)
	if err == nil {
		return price, nil
	}

	// Estrategia 3: Buscar usando regex en datos de mercado
	price, err = s.extractFromMarketData(htmlContent, symbol)
	if err == nil {
		return price, nil
	}

	return 0, fmt.Errorf("no se pudo extraer el precio para %s usando ninguna estrategia", symbol)
}

// extractFromFinStreamer busca el precio en elementos fin-streamer
func (s *AccionesStrategy) extractFromFinStreamer(htmlContent, symbol string) (float64, error) {
	// Patrón para buscar fin-streamer con el símbolo específico y regularMarketPrice
	pattern := fmt.Sprintf(`<fin-streamer[^>]*data-symbol="%s"[^>]*data-field="regularMarketPrice"[^>]*data-value="([0-9.]+)"[^>]*>`, strings.ToUpper(symbol))

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) >= 2 {
		price, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("error parseando precio desde fin-streamer: %w", err)
		}
		return price, nil
	}

	// Patrón alternativo buscando value= en cualquier posición
	altPattern := fmt.Sprintf(`data-symbol="%s"[^>]*data-field="regularMarketPrice"[^>]*value="([0-9.]+)"`, strings.ToUpper(symbol))
	altRe := regexp.MustCompile(altPattern)
	altMatches := altRe.FindStringSubmatch(htmlContent)

	if len(altMatches) >= 2 {
		price, err := strconv.ParseFloat(altMatches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("error parseando precio desde fin-streamer alternativo: %w", err)
		}
		return price, nil
	}

	return 0, fmt.Errorf("precio no encontrado en fin-streamer para símbolo %s", symbol)
}

// extractFromJSON busca el precio en objetos JSON embebidos
func (s *AccionesStrategy) extractFromJSON(htmlContent string) (float64, error) {
	// Buscar patrones de JSON con regularMarketPrice
	jsonPattern := `"regularMarketPrice":\s*{\s*"raw":\s*([0-9.]+)\s*,\s*"fmt":\s*"[^"]*"\s*}`

	re := regexp.MustCompile(jsonPattern)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) >= 2 {
		price, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("error parseando precio desde JSON: %w", err)
		}
		return price, nil
	}

	return 0, fmt.Errorf("precio no encontrado en JSON embebido")
}

// extractFromMarketData busca el precio en datos de mercado estructurados
func (s *AccionesStrategy) extractFromMarketData(htmlContent, symbol string) (float64, error) {
	// Buscar en contenido de texto de fin-streamer (fallback)
	pattern := fmt.Sprintf(`<fin-streamer[^>]*data-symbol="%s"[^>]*data-field="regularMarketPrice"[^>]*>([0-9,]+\.?[0-9]*)</fin-streamer>`, strings.ToUpper(symbol))

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) >= 2 {
		// Limpiar comas del precio
		priceStr := strings.ReplaceAll(matches[1], ",", "")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return 0, fmt.Errorf("error parseando precio desde contenido: %w", err)
		}
		return price, nil
	}

	return 0, fmt.Errorf("precio no encontrado en datos de mercado para símbolo %s", symbol)
}

// GetSupportedType devuelve el tipo soportado
func (s *AccionesStrategy) GetSupportedType() string {
	return "ACCIONES"
}

// BuildURL construye la URL usando la configuración del tipo de inversión
func (s *AccionesStrategy) BuildURL(baseURL, code string) string {
	// Limpiar espacios en blanco
	baseURL = strings.TrimSpace(baseURL)
	// Si la URL base ya contiene un placeholder {symbol} o {code}, reemplazarlo
	if strings.Contains(baseURL, "{symbol}") {
		return strings.ReplaceAll(baseURL, "{symbol}", code)
	}
	if strings.Contains(baseURL, "{code}") {
		return strings.ReplaceAll(baseURL, "{code}", code)
	}

	// Si la URL base termina con quote o quote/ asumir formato Yahoo Finance
	if strings.HasSuffix(baseURL, "quote/") {
		return fmt.Sprintf("%s%s/", baseURL, code)
	}
	if strings.HasSuffix(baseURL, "quote") {
		return fmt.Sprintf("%s/%s/", baseURL, code)
	}

	// Si la URL base termina con /, agregar el código
	if strings.HasSuffix(baseURL, "/") {
		return fmt.Sprintf("%s%s", baseURL, code)
	}

	// En otros casos, agregar el código con separador / y terminador /
	return fmt.Sprintf("%s/%s/", baseURL, code)
}
