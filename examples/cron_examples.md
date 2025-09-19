# Ejemplos de Uso del Servicio de Cron para Snapshots

Este documento proporciona ejemplos de cómo usar el servicio de cron para el scraping automático de assets y creación de snapshots.

## Descripción del Servicio

El servicio de cron se ejecuta automáticamente **todos los domingos a las 3:00 AM UTC** y realiza las siguientes tareas:

1. 🔍 **Obtiene todos los assets válidos** con sus tipos de inversión
2. 💰 **Scrapea el precio actual** de cada asset usando la estrategia correcta (Cedears, Criptomonedas, Acciones)
3. 📊 **Actualiza el campo `lastPrice`** del asset para optimizar futuras consultas
4. 📸 **Crea snapshots** para todos los holdings asociados a cada asset
5. 📈 **Calcula y actualiza las ganancias** (earnings) de cada holding

## Endpoints Disponibles

### 1. Obtener Estado del Cron

```bash
GET /api/admin/cron/status
Authorization: Bearer YOUR_API_KEY
```

**Respuesta:**

```json
{
  "status": "success",
  "message": "Estado del cron obtenido exitosamente",
  "data": {
    "running": true,
    "total_jobs": 1,
    "next_execution": "2024-03-17 03:00:00 UTC"
  }
}
```

### 2. Obtener Próxima Ejecución

```bash
GET /api/admin/cron/next
Authorization: Bearer YOUR_API_KEY
```

**Respuesta:**

```json
{
  "status": "success",
  "message": "Próxima ejecución obtenida exitosamente",
  "data": {
    "next_execution": "2024-03-17 03:00:00 UTC",
    "timezone": "UTC",
    "cron_expression": "0 3 * * 0",
    "description": "Se ejecuta todos los domingos a las 3:00 AM UTC"
  }
}
```

### 3. Ejecutar Scraping Manual

```bash
POST /api/admin/cron/execute
Authorization: Bearer YOUR_API_KEY
```

**Respuesta:**

```json
{
  "status": "success",
  "message": "Scraping manual iniciado en segundo plano",
  "data": {
    "message": "El scraping se está ejecutando. Revisa los logs del servidor para ver el progreso."
  }
}
```

### 4. Obtener Información General

```bash
GET /api/admin/cron/info
Authorization: Bearer YOUR_API_KEY
```

**Respuesta:**

```json
{
  "status": "success",
  "message": "Información del servicio de cron",
  "data": {
    "service_name": "Weekly Asset Scraping Service",
    "description": "Servicio que scrapea precios de assets y crea snapshots de holdings cada domingo",
    "schedule": "Domingos a las 3:00 AM UTC",
    "cron_expression": "0 3 * * 0",
    "status": {
      "running": true,
      "total_jobs": 1,
      "next_execution": "2024-03-17 03:00:00 UTC"
    },
    "next_execution": "2024-03-17 03:00:00 UTC",
    "features": [
      "Scraping automático de precios de assets",
      "Actualización del campo lastPrice para optimización",
      "Creación de snapshots para todos los holdings",
      "Cálculo automático de earnings",
      "Logging detallado del proceso",
      "Ejecución manual via API"
    ],
    "endpoints": [
      {
        "method": "GET",
        "path": "/api/admin/cron/status",
        "description": "Obtener estado del servicio de cron"
      },
      {
        "method": "GET",
        "path": "/api/admin/cron/next",
        "description": "Obtener próxima ejecución programada"
      },
      {
        "method": "POST",
        "path": "/api/admin/cron/execute",
        "description": "Ejecutar scraping manual"
      },
      {
        "method": "GET",
        "path": "/api/admin/cron/info",
        "description": "Obtener información general del servicio"
      }
    ]
  }
}
```

## Ejemplos con cURL

### Obtener estado del cron

```bash
curl -X GET "http://localhost:3000/api/admin/cron/status" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json"
```

### Ejecutar scraping manual

```bash
curl -X POST "http://localhost:3000/api/admin/cron/execute" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json"
```

### Obtener información completa

```bash
curl -X GET "http://localhost:3000/api/admin/cron/info" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json"
```

## Logs del Servidor

Cuando el cron se ejecuta, verás logs similares a estos en el servidor:

```
2024-03-17 03:00:00 🚀 Iniciando scraping semanal de assets...
2024-03-17 03:00:00 📊 Procesando 15 assets...
2024-03-17 03:00:01 🔍 Procesando asset: AAPL (AAPL)
2024-03-17 03:00:02 💰 Precio scrapeado para AAPL (AAPL): 150.25 USD
2024-03-17 03:00:02 📸 Creando 3 snapshots para asset AAPL (AAPL)
2024-03-17 03:00:02 📈 Earnings actualizados para holding abc-123: 25.50 (2.15%)
2024-03-17 03:00:02 ✅ Asset procesado exitosamente: AAPL (AAPL) - Precio: 150.25
...
2024-03-17 03:02:30 🏁 Scraping semanal completado en 2m30s - Éxitos: 14, Errores: 1
```

## Optimizaciones Implementadas

### 1. **Campo `lastPrice` en Assets**

- Se actualiza con cada scraping para evitar múltiples consultas
- Permite obtener el último precio conocido sin hacer scraping

### 2. **Pausa entre Requests**

- 200ms entre cada asset para ser respetuosos con los servidores
- Evita ser bloqueado por rate limiting

### 3. **Manejo de Errores**

- Continúa procesando otros assets si uno falla
- Logging detallado de errores y éxitos

### 4. **Cálculo Automático de Earnings**

- Compara con el snapshot anterior para calcular ganancias
- Actualiza tanto earnings absolutos como relativos (%)

## Consideraciones de Horario

- **Horario programado**: Domingos 3:00 AM UTC
- **Razón**: Minimiza impacto en horarios de trading
- **Duración estimada**: 2-5 minutos dependiendo del número de assets

## Monitoreo y Troubleshooting

### Verificar si el cron está funcionando:

```bash
curl -X GET "http://localhost:3000/api/admin/cron/status" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### Ejecutar manualmente para testing:

```bash
curl -X POST "http://localhost:3000/api/admin/cron/execute" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### Revisar logs del servidor:

```bash
docker-compose logs -f app
```

## Arquitectura del Sistema

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CronService   │───▶│ ScrapingService │───▶│   Strategies    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Assets      │    │    Holdings     │    │   Yahoo Finance │
│   (lastPrice)   │    │   (earnings)    │    │   Other Sources │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│   Snapshots     │    │   Database      │
│   (historical)  │    │   (persistent)  │
└─────────────────┘    └─────────────────┘
```
