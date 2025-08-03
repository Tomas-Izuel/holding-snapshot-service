# 🔌 Ejemplos de Uso de la API

## Configuración Inicial

```bash
# Configurar API Key
export API_KEY="tu_api_key_aqui"
```

## Ejemplos de Requests

### 1. Health Check

```bash
curl -X GET http://localhost:8080/api/health
```

**Respuesta:**

```json
{
  "status": "ok",
  "service": "holding-snapshots",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 2. Validar Holding - Caso Exitoso

```bash
curl -X POST http://localhost:8080/api/validate \
  -H "Content-Type: application/json" \
  -H "Authorization: $API_KEY" \
  -d '{
    "name": "Apple Inc.",
    "code": "AAPL",
    "quantity": 10.5,
    "groupId": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

**Respuesta:**

```json
{
  "holding": {
    "name": "Apple Inc.",
    "code": "AAPL",
    "quantity": 10.5,
    "group_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "is_valid": true
}
```

### 3. Validar Holding - Caso de Error

```bash
curl -X POST http://localhost:8080/api/validate \
  -H "Content-Type: application/json" \
  -H "Authorization: $API_KEY" \
  -d '{
    "name": "Activo Inexistente",
    "code": "FAKE",
    "quantity": 5,
    "groupId": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

**Respuesta:**

```json
{
  "holding": {
    "name": "Activo Inexistente",
    "code": "FAKE",
    "quantity": 5,
    "group_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "is_valid": false
}
```

### 4. Ejecutar Scraping Manual

```bash
curl -X POST http://localhost:8080/api/scraping/execute \
  -H "Authorization: $API_KEY"
```

**Respuesta:**

```json
{
  "message": "Scraping manual iniciado"
}
```

## Códigos de Error

| Código | Descripción                           |
| ------ | ------------------------------------- |
| 200    | Solicitud exitosa                     |
| 400    | Request mal formado o datos inválidos |
| 401    | API Key faltante o inválida           |
| 500    | Error interno del servidor            |

## Notas Importantes

1. **API Key**: Siempre incluir la API Key en el header `Authorization`
2. **Content-Type**: Usar `application/json` para requests POST
3. **UUIDs**: Los `groupId` deben ser UUIDs válidos
4. **Quantity**: Debe ser un número mayor a 0
5. **Rate Limiting**: No implementado, pero se recomienda no exceder 10 requests/segundo
