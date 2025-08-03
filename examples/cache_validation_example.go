package main

import (
	"fmt"
	"holding-snapshots/internal/services"
	"log"
)

// Ejemplo de uso del cache de validación de holdings
func mainCacheValidation() {
	fmt.Println("🚀 Ejemplo del Cache de Validación de Holdings")
	fmt.Println("=" + repeatStringCache("=", 45))

	// Crear el servicio de scraping
	scrapingService := services.NewScrapingService()

	// Simular IDs de grupos y tipos de inversión
	groupCedears := "grupo-cedears-uuid"
	groupCrypto := "grupo-crypto-uuid"
	groupAcciones := "grupo-acciones-uuid"

	fmt.Println("\n📋 Ejemplos de validación con cache:")
	fmt.Println(repeatStringCache("-", 40))

	// Ejemplo 1: Primera validación (se hace scraping)
	fmt.Println("\n🔍 PRIMERA VALIDACIÓN - CRYPTO")
	holding1, valid1, err1 := scrapingService.ValidateHolding(
		"Bitcoin",   // name
		"BTC",       // code
		groupCrypto, // groupID
		1.5,         // quantity
	)

	if err1 != nil {
		log.Printf("❌ Error: %v", err1)
	} else if valid1 {
		fmt.Printf("✅ Bitcoin validado exitosamente: %+v\n", holding1.Name)
		fmt.Println("   💾 Resultado guardado en cache por 24 horas")
	} else {
		fmt.Println("❌ Bitcoin no es válido")
	}

	// Ejemplo 2: Segunda validación del mismo activo (se usa cache)
	fmt.Println("\n📦 SEGUNDA VALIDACIÓN DEL MISMO ACTIVO - CRYPTO")
	holding2, valid2, err2 := scrapingService.ValidateHolding(
		"Bitcoin",   // mismo name
		"BTC",       // mismo code
		groupCrypto, // mismo groupID
		2.0,         // diferente quantity
	)

	if err2 != nil {
		log.Printf("❌ Error: %v", err2)
	} else if valid2 {
		fmt.Printf("✅ Bitcoin validado desde cache: %+v\n", holding2.Name)
		fmt.Println("   ⚡ No se hizo scraping, se usó cache!")
	}

	// Ejemplo 3: Validación de activo diferente (se hace scraping)
	fmt.Println("\n🔍 VALIDACIÓN DE ACTIVO DIFERENTE - CEDEARS")
	holding3, valid3, err3 := scrapingService.ValidateHolding(
		"Apple Inc",  // name
		"AAPL",       // code
		groupCedears, // groupID diferente
		10.0,         // quantity
	)

	if err3 != nil {
		log.Printf("❌ Error: %v", err3)
	} else if valid3 {
		fmt.Printf("✅ AAPL validado exitosamente: %+v\n", holding3.Name)
		fmt.Println("   💾 Resultado guardado en cache por 24 horas")
	}

	// Ejemplo 4: Validación del mismo código pero diferente tipo
	fmt.Println("\n🔍 MISMO CÓDIGO, DIFERENTE TIPO - ACCIONES")
	fmt.Println("   (AAPL como acción vs AAPL como CEDEAR)")
	holding4, valid4, err4 := scrapingService.ValidateHolding(
		"Apple Inc Stock", // name
		"AAPL",            // mismo code
		groupAcciones,     // tipo diferente
		5.0,               // quantity
	)

	if err4 != nil {
		log.Printf("❌ Error: %v", err4)
	} else if valid4 {
		fmt.Printf("✅ AAPL como acción validado: %+v\n", holding4.Name)
		fmt.Println("   💾 Cache separado por tipo de inversión")
	}

	fmt.Println("\n📊 Información del Cache:")
	fmt.Println(repeatStringCache("-", 30))

	// Obtener estadísticas del cache
	stats := scrapingService.GetValidationCacheStats()
	for key, value := range stats {
		fmt.Printf("• %s: %v\n", key, value)
	}

	fmt.Println("\n🔧 Gestión del Cache:")
	fmt.Println(repeatStringCache("-", 25))

	// Ejemplo de obtener datos desde cache
	fmt.Println("\n📦 Obtener datos desde cache:")
	fmt.Println("// Obtener BTC desde cache")
	fmt.Println("cachedData, found, err := scrapingService.GetValidatedHoldingFromCache(\"crypto-type-id\", \"BTC\")")
	fmt.Println("if found {")
	fmt.Println("    fmt.Printf(\"Datos en cache: %+v\", cachedData)")
	fmt.Println("}")

	// Ejemplo de limpiar cache
	fmt.Println("\n🗑️ Limpiar cache específico:")
	fmt.Println("// Limpiar cache de BTC")
	fmt.Println("err := scrapingService.ClearValidatedHoldingCache(\"crypto-type-id\", \"BTC\")")
	fmt.Println("if err != nil {")
	fmt.Println("    log.Printf(\"Error limpiando cache: %v\", err)")
	fmt.Println("}")

	fmt.Println("\n💡 Ventajas del Cache de Validación:")
	fmt.Println(repeatStringCache("-", 40))
	fmt.Println("• ⚡ Evita scraping repetido para activos ya validados")
	fmt.Println("• 🎯 Cache separado por tipo de inversión (CRYPTO, CEDEARS, ACCIONES)")
	fmt.Println("• ⏰ TTL diferenciado: 24h válidos, 2h inválidos")
	fmt.Println("• 💾 Persistencia en Redis para múltiples instancias")
	fmt.Println("• 📊 Información detallada sobre validaciones previas")
	fmt.Println("• 🔧 Métodos para gestionar el cache manualmente")

	fmt.Println("\n📝 Estructura de Clave de Cache:")
	fmt.Println(repeatStringCache("-", 35))
	fmt.Println("Patrón: validated_holding:{typeID}:{code}")
	fmt.Println("Ejemplos:")
	fmt.Println("• validated_holding:crypto-uuid:BTC")
	fmt.Println("• validated_holding:cedears-uuid:AAPL")
	fmt.Println("• validated_holding:acciones-uuid:AAPL")

	fmt.Println("\n🔄 Flujo de Validación:")
	fmt.Println(repeatStringCache("-", 25))
	fmt.Println("1. 🔍 Buscar en cache Redis")
	fmt.Println("2. 📦 Si existe → Devolver resultado cached")
	fmt.Println("3. 🌐 Si no existe → Hacer scraping")
	fmt.Println("4. 💾 Guardar resultado en cache con TTL")
	fmt.Println("5. ✅ Devolver resultado al usuario")
}

// Función helper para crear líneas de separación (cache validation)
func repeatStringCache(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
