package main

import (
	"fmt"
	"holding-snapshots/internal/scraping"
)

// Ejemplo de uso del sistema de scraping con Factory Pattern
func main() {
	fmt.Println("🏭 Ejemplo de uso del Scraping Factory Pattern")
	fmt.Println(repeatString("=", 50))

	// Crear la factory
	factory := scraping.NewScrapingFactory()

	// Ejemplos de diferentes tipos de grupos
	ejemplosGrupos := []struct {
		nombre      string
		descripcion string
	}{
		{"MIS CEDEARS", "Grupo de CEDEARs argentinos"},
		{"ACCIONES USA", "Acciones estadounidenses"},
		{"CRYPTO PORTFOLIO", "Cartera de criptomonedas"},
		{"MIS CRIPTOS", "Otra forma de nombrar crypto"},
		{"BONOS ARGENTINA", "Tipo no soportado aún"},
	}

	fmt.Println("\n📋 Validación de nombres de grupos:")
	fmt.Println(repeatString("-", 40))

	for _, ejemplo := range ejemplosGrupos {
		valido, tipoEstrategia, mensaje := factory.ValidateGroupName(ejemplo.nombre)
		status := "✅"
		if !valido {
			status = "❌"
		}

		fmt.Printf("%s %s (%s)\n", status, ejemplo.nombre, ejemplo.descripcion)
		fmt.Printf("   Estrategia: %s\n", tipoEstrategia)
		fmt.Printf("   Mensaje: %s\n\n", mensaje)
	}

	fmt.Println("\n🔍 Estrategias disponibles:")
	fmt.Println(repeatString("-", 30))
	estrategias := factory.GetAvailableStrategies()
	for _, estrategia := range estrategias {
		fmt.Printf("• %s\n", estrategia)
	}

	fmt.Println("\n🧪 Simulación de obtención de estrategias:")
	fmt.Println(repeatString("-", 45))

	// Simular obtención de estrategias para diferentes grupos
	gruposEjemplo := []string{"MIS CEDEARS", "ACCIONES TECH", "CRYPTO WALLET"}

	for _, grupoNombre := range gruposEjemplo {
		fmt.Printf("\n📊 Procesando grupo: %s\n", grupoNombre)

		if strategy, err := factory.GetStrategy(grupoNombre); err == nil {
			fmt.Printf("   ✅ Estrategia obtenida: %s\n", strategy.GetSupportedType())

			// Ejemplo de construcción de URL
			exampleURL := strategy.BuildURL("https://api.ejemplo.com/price", "AAPL")
			fmt.Printf("   🔗 URL de ejemplo: %s\n", exampleURL)
		} else {
			fmt.Printf("   ❌ Error: %v\n", err)
		}
	}

	fmt.Println("\n💡 Ejemplo de uso en ScrapingService:")
	fmt.Println(repeatString("-", 42))
	fmt.Println()
	fmt.Println("// En tu código:")
	fmt.Println("scrapingService := services.NewScrapingService()")
	fmt.Println()
	fmt.Println("// Ejemplo para CEDEARS")
	fmt.Println("typeInvestment := &models.TypeInvestment{")
	fmt.Println("    Name:        \"Cedears\",")
	fmt.Println("    ScrapingURL: \"https://api.cedears.com/price\",")
	fmt.Println("    Currency:    \"ARS\",")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("price, err := scrapingService.FetchAssetPrice(")
	fmt.Println("    typeInvestment,")
	fmt.Println("    \"AAPL\",")
	fmt.Println("    \"MIS CEDEARS\"  // La factory detecta automáticamente que debe usar CedearsStrategy")
	fmt.Println(")")
	fmt.Println()
	fmt.Println("if err != nil {")
	fmt.Println("    log.Printf(\"Error: %%v\", err)")
	fmt.Println("} else {")
	fmt.Println("    fmt.Printf(\"Precio de AAPL: %%.2f ARS\", price)")
	fmt.Println("}")

	fmt.Println("\n🎯 Ventajas del Factory Pattern:")
	fmt.Println(repeatString("-", 35))
	fmt.Println("• ✨ Código más limpio y modular")
	fmt.Println("• 🔧 Fácil mantenimiento de cada estrategia")
	fmt.Println("• 📈 Extensible para nuevos tipos de activos")
	fmt.Println("• ⚡ Configuraciones específicas por tipo")
	fmt.Println("• 🧪 Testing independiente por estrategia")
	fmt.Println("• 🎛️ Control granular de timeouts y cache")
}

// Función helper para crear líneas de separación
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
