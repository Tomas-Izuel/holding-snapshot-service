package scraping

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"holding-snapshots/internal/models"

	"github.com/gocolly/colly"
)

type StockStrategy struct{}

// FetchPrice obtiene el precio de una acción usando web scraping de Yahoo Finance
func (s *StockStrategy) FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error) {
	log.Printf("🔍 [StockStrategy] FetchPrice iniciado - TypeInvestment: %s, Code: %s", typeInvestment.Name, code)
	log.Printf("🔍 [StockStrategy] ScrapingURL base: %s", typeInvestment.ScrapingURL)

	// Construir la URL específica para el código de la acción
	url := s.BuildURL(typeInvestment.ScrapingURL, code)

	log.Printf("🌐 [StockStrategy] URL construida: %s", url)
	// Instanciar un nuevo colector
	c := colly.NewCollector(
		colly.AllowedDomains("finance.yahoo.com"),
	)

	log.Printf("🤖 [StockStrategy] Colector creado correctamente")

	var price string
	var found bool

	// OnHTML callback para extraer información de Yahoo Finance
	c.OnHTML("section[data-testid=\"quote-hdr\"]", func(e *colly.HTMLElement) {
		log.Printf("📊 [StockStrategy] Elemento quote-hdr encontrado")
		// Extraer el precio del elemento correspondiente
		priceText := e.ChildText("[data-testid=\"qsp-price\"]")
		log.Printf("💰 [StockStrategy] Precio extraído: '%s'", priceText)
		if priceText != "" {
			price = priceText
			found = true
			log.Printf("✅ [StockStrategy] Precio válido encontrado: %s", price)
		}
	})

	// Manejar errores durante el scraping
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("❌ [StockStrategy] Error durante el scraping de %s: %v", url, err)
	})

	// Callback para verificar la respuesta HTTP
	c.OnResponse(func(r *colly.Response) {
		log.Printf("📡 [StockStrategy] Respuesta HTTP recibida - Status: %d, URL: %s", r.StatusCode, r.Request.URL.String())
	})

	// Realizar la solicitud HTTP
	log.Printf("🚀 [StockStrategy] Iniciando visita a URL: %s", url)
	err := c.Visit(url)
	if err != nil {
		log.Printf("❌ [StockStrategy] Error al visitar la URL: %v", err)
		return 0, fmt.Errorf("error al visitar la URL %s: %v", url, err)
	}

	log.Printf("🔍 [StockStrategy] Visita completada - Found: %t, Price: '%s'", found, price)

	// Verificar si se encontró el precio
	if !found || price == "" {
		log.Printf("⚠️ [StockStrategy] No se encontró precio válido")
		return 0, fmt.Errorf("no se pudo encontrar el precio para el código %s en la URL %s", code, url)
	}

	// Limpiar y convertir el precio a float64
	cleanPrice := strings.ReplaceAll(price, ",", "")
	priceFloat, err := strconv.ParseFloat(cleanPrice, 64)
	if err != nil {
		return 0, fmt.Errorf("error al convertir el precio '%s' a número: %v", price, err)
	}

	return priceFloat, nil
}

// BuildURL construye la URL específica para el scraping de Yahoo Finance
func (s *StockStrategy) BuildURL(baseURL, code string) string {
	// Usar el ScrapingURL proporcionado y agregar el código de la acción
	// Para acciones: https://finance.yahoo.com/quote/SPY
	return fmt.Sprintf("%s/%s", baseURL, code)
}
