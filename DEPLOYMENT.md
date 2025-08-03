# 🚀 Guía de Deployment - Holding Snapshots Service

Esta guía te ayudará a desplegar el servicio de snapshots como un microservicio independiente que se conecta al servicio principal.

## 📋 Prerequisitos

Antes de desplegar este servicio, asegúrate de que:

1. **El servicio principal esté ejecutándose** y accessible
2. **PostgreSQL esté disponible** con las tablas del sistema principal
3. **Redis esté disponible** para cache
4. Tengas las **credenciales de conexión** a las bases de datos
5. **API Key configurada** para la comunicación entre servicios

## 🔧 Configuración de Variables de Entorno

### Variables Requeridas

```bash
# Conexión a la base de datos del servicio principal
DATABASE_URL=postgresql://usuario:password@host:puerto/nombre_db

# Conexión al Redis del servicio principal
REDIS_URL=redis://host:puerto

# API Key para autenticación (debe coincidir con el servicio principal)
SNAPSHOT_SERVICE_API_KEY=tu_api_key_secreta

# Configuración del servidor
PORT=8080
ENV=production

# Configuración del cron job (domingos 1:00 AM)
SCRAPING_CRON_SCHEDULE=0 1 * * 0
```

### Ejemplo de configuración para desarrollo local

```bash
# Si el servicio principal está en localhost
DATABASE_URL=postgresql://holding_admin:holding_password@localhost:5432/holdingdb
REDIS_URL=redis://localhost:6379
SNAPSHOT_SERVICE_API_KEY=dev_api_key_12345
```

### Ejemplo de configuración para producción

```bash
# Si el servicio principal está en un servidor remoto
DATABASE_URL=postgresql://holding_user:secure_password@10.0.1.100:5432/holding_production
REDIS_URL=redis://10.0.1.100:6379
SNAPSHOT_SERVICE_API_KEY=prod_super_secure_api_key_2024
```

## 🐳 Deployment con Docker

### 1. Preparar el archivo .env

```bash
# Crear archivo .env con las configuraciones de tu entorno
cat > .env << 'EOF'
DATABASE_URL=postgresql://tu_usuario:tu_password@tu_host:5432/tu_db
REDIS_URL=redis://tu_host:6379
SNAPSHOT_SERVICE_API_KEY=tu_api_key_real
PORT=8080
ENV=production
SCRAPING_CRON_SCHEDULE=0 1 * * 0
EOF
```

### 2. Construir y ejecutar con Docker Compose

```bash
# Construir la imagen y ejecutar el servicio
docker-compose up --build -d

# Verificar que el servicio esté ejecutándose
docker-compose ps

# Ver logs del servicio
docker-compose logs -f holding-snapshots
```

### 3. Verificar la conexión

```bash
# Test de salud del servicio
curl http://localhost:8080/api/health

# Debería retornar:
# {"status":"ok","service":"holding-snapshots","timestamp":"..."}
```

## 🖥️ Deployment sin Docker

### 1. Compilar la aplicación

```bash
# Instalar dependencias
go mod download

# Compilar para producción
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o holding-snapshots cmd/server/main.go
```

### 2. Configurar variables de entorno

```bash
# Exportar variables de entorno
export DATABASE_URL="postgresql://usuario:password@host:port/db"
export REDIS_URL="redis://host:port"
export SNAPSHOT_SERVICE_API_KEY="tu_api_key"
export PORT=8080
export ENV=production
```

### 3. Ejecutar el servicio

```bash
# Ejecutar la aplicación
./holding-snapshots

# O como servicio de sistema (systemd)
sudo systemctl start holding-snapshots
```

## ⚙️ Configuración como Servicio de Sistema (Linux)

### 1. Crear archivo de servicio systemd

```bash
sudo nano /etc/systemd/system/holding-snapshots.service
```

```ini
[Unit]
Description=Holding Snapshots Service
After=network.target

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/holding-snapshots
ExecStart=/opt/holding-snapshots/holding-snapshots
Restart=always
RestartSec=5

# Variables de entorno
Environment=DATABASE_URL=postgresql://usuario:password@host:port/db
Environment=REDIS_URL=redis://host:port
Environment=SNAPSHOT_SERVICE_API_KEY=tu_api_key
Environment=PORT=8080
Environment=ENV=production

[Install]
WantedBy=multi-user.target
```

### 2. Habilitar y iniciar el servicio

```bash
# Recargar configuración de systemd
sudo systemctl daemon-reload

# Habilitar el servicio para que inicie automáticamente
sudo systemctl enable holding-snapshots

# Iniciar el servicio
sudo systemctl start holding-snapshots

# Verificar estado
sudo systemctl status holding-snapshots
```

## 🔍 Verificación del Deployment

### 1. Health Check

```bash
curl http://localhost:8080/api/health
```

### 2. Test de validación (requiere API Key)

```bash
curl -X POST http://localhost:8080/api/validate \
  -H "Content-Type: application/json" \
  -H "Authorization: tu_api_key" \
  -d '{
    "name": "Test Asset",
    "code": "TEST",
    "quantity": 1,
    "groupId": "uuid-valido-del-grupo"
  }'
```

### 3. Verificar logs

```bash
# Docker Compose
docker-compose logs -f

# Systemd
sudo journalctl -f -u holding-snapshots

# Archivo de logs (si se configura)
tail -f /var/log/holding-snapshots.log
```

## 🚨 Troubleshooting

### Error de conexión a base de datos

```bash
# Verificar conectividad
telnet host_db 5432

# Verificar credenciales
psql -h host_db -p 5432 -U usuario -d nombre_db
```

### Error de conexión a Redis

```bash
# Verificar conectividad
telnet host_redis 6379

# Test con redis-cli
redis-cli -h host_redis -p 6379 ping
```

### Error de autenticación API

- Verificar que la `SNAPSHOT_SERVICE_API_KEY` sea exactamente la misma en ambos servicios
- Verificar que el header `Authorization` se esté enviando correctamente

### El cron job no se ejecuta

- Verificar que el formato del cron schedule sea correcto
- Verificar logs para errores durante la ejecución del scraping
- Verificar que el servicio tenga acceso a las URLs de scraping

## 📊 Monitoreo

### Endpoints de monitoreo

- **Health Check**: `GET /api/health`
- **Métricas**: Los logs proporcionan información detallada sobre:
  - Conexiones a base de datos
  - Ejecución de cron jobs
  - Requests de validación
  - Errores de scraping

### Logs importantes a monitorear

```bash
# Conexión exitosa a DB
"✅ Conexión a la base de datos establecida"

# Conexión exitosa a Redis
"✅ Conexión a Redis establecida"

# Inicio de cron job
"🚀 Iniciando scraping semanal..."

# Scraping exitoso
"✅ Scraping semanal completado"
```

## 🔄 Actualizaciones

### Para actualizar el servicio:

```bash
# 1. Detener el servicio
docker-compose down

# 2. Actualizar código
git pull origin main

# 3. Reconstruir y reiniciar
docker-compose up --build -d

# 4. Verificar que funcione
curl http://localhost:8080/api/health
```
