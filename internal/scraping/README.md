# Sistema de Scraping con Factory Pattern

Este paquete implementa un sistema de scraping modular usando el patrón Factory para manejar diferentes tipos de activos financieros.

## Arquitectura

### 1. Interface ScrapingStrategy

Define el contrato que deben cumplir todas las estrategias de scraping:

- `FetchPrice()`: Obtiene el precio del activo
- `GetSupportedType()`: Devuelve el tipo soportado
- `BuildURL()`: Construye la URL específica para el scraping

### 2. Implementaciones Concretas

#### CedearsStrategy

- **Tipo soportado**: `CEDEARS`
- **Características**:
  - Cache de 10 minutos
  - Headers específicos para CEDEARs
  - Parámetros: `market=cedears&currency=ars`
  - Timeout: 15 segundos

#### AccionesStrategy

- **Tipo soportado**: `ACCIONES`
- **Características**:
  - Cache de 5 minutos (más frecuente)
  - Headers de navegador estándar
  - Parámetros: `type=stock&exchange=NYSE,NASDAQ`
  - Timeout: 12 segundos

#### CryptoStrategy

- **Tipo soportado**: `CRYPTO`
- **Características**:
  - Cache de 3 minutos (muy frecuente)
  - Headers específicos para APIs crypto
  - Parámetros: `convert=USD,ARS&include_market_data=true`
  - Manejo de precios USD/ARS
  - Timeout: 10 segundos

### 3. ScrapingFactory

La factory mapea nombres de grupos a estrategias:

- **CEDEARS**: Grupos que contengan "CEDEAR"
- **ACCIONES**: Grupos que contengan "ACCION"
- **CRYPTO**: Grupos que contengan "CRYPTO" o "CRIPTO"

## Uso

### Scraping de Precios

```go
// Crear el servicio de scraping
scrapingService := NewScrapingService()

// El servicio automáticamente usa la factory para obtener la estrategia correcta
price, err := scrapingService.FetchAssetPrice(&typeInvestment, "AAPL", "MIS ACCIONES")

// La factory detecta que "MIS ACCIONES" contiene "ACCION" y usa AccionesStrategy
```

### Validación con Cache de Holdings

```go
// Primera validación (hace scraping y guarda en cache)
holding, valid, err := scrapingService.ValidateHolding("Bitcoin", "BTC", groupID, 1.5)

// Segunda validación del mismo activo (usa cache, no hace scraping)
holding2, valid2, err2 := scrapingService.ValidateHolding("Bitcoin", "BTC", groupID, 2.0)

// Gestión manual del cache
cachedData, found, err := scrapingService.GetValidatedHoldingFromCache(typeID, "BTC")
err = scrapingService.ClearValidatedHoldingCache(typeID, "BTC")
```

## Ventajas

1. **Extensibilidad**: Fácil agregar nuevos tipos de activos
2. **Mantenimiento**: Cada estrategia es independiente
3. **Configuración**: Cada tipo tiene su propia configuración de timeout, cache, headers
4. **Testing**: Se puede testear cada estrategia por separado
5. **Flexibilidad**: URLs y parámetros específicos por tipo
6. **Cache Inteligente**: Evita scraping repetido con cache Redis
7. **Separación por Tipos**: Cache independiente por tipo de inversión
8. **TTL Optimizado**: 24h para válidos, 2h para inválidos

## Agregando Nueva Estrategia

1. Crear archivo `nuevo_tipo_strategy.go`
2. Implementar interface `ScrapingStrategy`
3. Registrar en `NewScrapingFactory()`
4. Agregar mapeo en `GetStrategy()`

```go
// Ejemplo: BondStrategy
type BondStrategy struct{}

func (s *BondStrategy) FetchPrice(typeInvestment *models.TypeInvestment, code string) (float64, error) {
    // Implementación específica para bonos
}

func (s *BondStrategy) GetSupportedType() string {
    return "BONOS"
}

func (s *BondStrategy) BuildURL(baseURL, code string) string {
    return fmt.Sprintf("%s?symbol=%s&type=bond&market=NYSE", baseURL, code)
}
```

## Cache de Validación de Holdings

### Estructura del Cache

- **Patrón de clave**: `validated_holding:{typeID}:{code}`
- **TTL para válidos**: 24 horas
- **TTL para inválidos**: 2 horas
- **Storage**: Redis

### Beneficios del Cache

1. **Performance**: Evita scraping repetido para activos ya validados
2. **Separación**: Cache independiente por tipo de inversión
3. **Persistencia**: Compartido entre múltiples instancias de la aplicación
4. **Gestión**: Métodos para administrar el cache manualmente

### Ejemplos de Claves de Cache

```
validated_holding:crypto-uuid-123:BTC
validated_holding:cedears-uuid-456:AAPL
validated_holding:acciones-uuid-789:AAPL
```

### Métodos de Gestión

```go
// Obtener desde cache
cachedData, found, err := service.GetValidatedHoldingFromCache(typeID, code)

// Limpiar cache específico
err := service.ClearValidatedHoldingCache(typeID, code)

// Obtener estadísticas
stats := service.GetValidationCacheStats()
```

### Flujo de Validación

1. 🔍 Buscar en cache Redis con clave `validated_holding:{typeID}:{code}`
2. 📦 Si existe y es válido → Devolver resultado cached
3. 🌐 Si no existe → Ejecutar scraping usando estrategia apropiada
4. 💾 Guardar resultado en cache con TTL correspondiente
5. ✅ Devolver resultado final al usuario
