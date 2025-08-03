# 📄 Holding Snapshots Service

Microservicio de scraping semanal para activos de inversión construido en Go con arquitectura MVC clásica. Se conecta al servicio principal para realizar snapshots periódicos de precios de activos.

## 🚀 Características

- **Microservicio independiente**: Se conecta a las bases de datos del servicio principal
- **Scraping automático**: Cron job que ejecuta scraping semanal los domingos a la 1:00 AM
- **Validación de holdings**: Endpoint para validar activos antes de agregarlos
- **Cache inteligente**: Utiliza Redis del servicio principal para optimizar requests
- **API REST**: Endpoints mínimos para comunicación con el servicio principal
- **Dockerizado**: Configuración simplificada para deployment como microservicio
- **Sin dependencias de DB**: Se conecta a PostgreSQL y Redis existentes

## 📦 Stack Tecnológico

- **Lenguaje**: Go 1.21
- **Framework HTTP**: Fiber v2
- **ORM**: GORM
- **Base de datos**: PostgreSQL 15
- **Cache**: Redis 7
- **Containerización**: Docker & Docker Compose

## 🏗️ Arquitectura

```
holding-snapshots/
├── cmd/server/           # Punto de entrada de la aplicación
├── internal/
│   ├── config/          # Configuración de la aplicación
│   ├── controllers/     # Controladores HTTP (Capa de presentación)
│   ├── models/          # Modelos de datos (Capa de modelo)
│   ├── services/        # Lógica de negocio (Capa de servicio)
│   ├── middleware/      # Middlewares HTTP
│   └── routes/          # Definición de rutas
├── pkg/
│   ├── database/        # Configuración de base de datos
│   └── cache/           # Configuración de cache
├── docker/              # Archivos de configuración Docker
└── migrations/          # Migraciones SQL (si necesario)
```

## 🗃️ Modelos de Datos

- **User**: Usuarios del sistema
- **TypeUser**: Tipos de usuario con permisos
- **Permission**: Permisos específicos
- **Group**: Grupos de inversión de usuarios
- **TypeInvestment**: Tipos de inversión con URLs de scraping
- **Holding**: Activos de inversión específicos
- **Snapshot**: Snapshots de precios en momentos específicos

## 🔌 Endpoints API

### Validar Holding

```http
POST /api/validate
Authorization: YOUR_API_KEY
Content-Type: application/json

{
  "name": "Apple Inc.",
  "code": "AAPL",
  "quantity": 10.5,
  "groupId": "uuid-del-grupo"
}
```

**Respuesta:**

```json
{
  "holding": {
    "name": "Apple Inc.",
    "code": "AAPL",
    "quantity": 10.5,
    "group_id": "uuid-del-grupo"
  },
  "is_valid": true
}
```

### Health Check

```http
GET /api/health
```

### Scraping Manual (para testing)

```http
POST /api/scraping/execute
Authorization: YOUR_API_KEY
```

## 🚀 Instalación y Configuración

### Prerequisitos

- Go 1.21 o superior
- Docker y Docker Compose
- Make (opcional, para usar Makefile)
- **Servicio principal ejecutándose** con PostgreSQL y Redis

### 1. Clonar el repositorio

```bash
git clone <repository-url>
cd holding-snapshots
```

### 2. Configurar variables de entorno

```bash
# Crear archivo .env con las configuraciones de conexión al servicio principal
cat > .env << EOF
DATABASE_URL=postgresql://usuario:password@host:puerto/dbname
REDIS_URL=redis://host:puerto
SNAPSHOT_SERVICE_API_KEY=tu_api_key_secreta
PORT=8080
ENV=development
SCRAPING_CRON_SCHEDULE=0 1 * * 0
EOF

# Editar variables según tu configuración del servicio principal
nano .env
```

### 3. Usando Docker Compose (Recomendado)

```bash
# Construir y ejecutar el servicio
docker-compose up --build -d

# Ver logs
docker-compose logs -f
```

### 4. Desarrollo local

```bash
# Instalar dependencias
make deps

# Ejecutar la aplicación
make run
```

⚠️ **Importante**: Este servicio requiere que el servicio principal esté ejecutándose con PostgreSQL y Redis disponibles.

## 📊 Base de Datos

El servicio se conecta a la base de datos PostgreSQL del servicio principal, que debe tener:

- Extensión `uuid-ossp` habilitada
- Esquema de tablas del sistema principal
- Permisos adecuados para el usuario de conexión

### Migraciones

Las migraciones se ejecutan automáticamente al iniciar la aplicación usando GORM AutoMigrate. Solo se crearán las tablas que no existan previamente.

## ⏰ Cron Jobs

El servicio incluye un cron job configurado para ejecutarse:

- **Scraping semanal**: Domingos a la 1:00 AM (configurable via `SCRAPING_CRON_SCHEDULE`)

## 🛠️ Comandos Útiles

```bash
# Compilar la aplicación
make build

# Ejecutar tests
make test

# Limpiar archivos compilados
make clean

# Ejecutar linters
make lint

# Parar servicios Docker
make docker-down

# Entrar al contenedor de la aplicación
make shell
```

## 🔧 Configuración

### Variables de Entorno

| Variable                   | Descripción                  | Valor por defecto        |
| -------------------------- | ---------------------------- | ------------------------ |
| `DATABASE_URL`             | URL de conexión a PostgreSQL | -                        |
| `REDIS_URL`                | URL de conexión a Redis      | `redis://localhost:6379` |
| `PORT`                     | Puerto del servidor          | `8080`                   |
| `ENV`                      | Entorno de ejecución         | `development`            |
| `SNAPSHOT_SERVICE_API_KEY` | API Key para autenticación   | -                        |
| `SCRAPING_CRON_SCHEDULE`   | Schedule del cron job        | `0 1 * * 0`              |

## 🧠 Comportamiento del Servicio

### Scraping Automático

1. **Trigger**: Cron job los domingos a la 1:00 AM
2. **Proceso**:
   - Obtiene todos los holdings de la base de datos
   - Para cada holding, consulta el precio actual en la URL de scraping correspondiente
   - Crea un nuevo snapshot con el precio actual
   - Actualiza el holding con earnings y precio actual
   - Calcula ganancias relativas comparando con snapshots anteriores

### Validación de Holdings

1. **Trigger**: Request POST a `/api/validate`
2. **Proceso**:
   - Valida el formato de la request
   - Consulta la URL de scraping del tipo de inversión
   - Verifica si el activo existe y es válido
   - Retorna resultado de validación

## 🔒 Seguridad

- **Autenticación**: API Key en header Authorization
- **CORS**: Configurado para permitir requests del frontend
- **Input Validation**: Validación de todos los inputs de API
- **Error Handling**: Manejo robusto de errores

## 📈 Monitoreo

- **Health Check**: Endpoint `/api/health` para verificar estado del servicio
- **Logs estructurados**: Logs detallados de todas las operaciones
- **Adminer**: Interface web para administrar PostgreSQL (puerto 8081)

## 🚧 Desarrollo

### Agregar nuevos endpoints

1. Crear controlador en `internal/controllers/`
2. Implementar lógica de negocio en `internal/services/`
3. Registrar rutas en `internal/routes/routes.go`

### Agregar nuevos modelos

1. Crear modelo en `internal/models/`
2. Agregar al AutoMigrate en `pkg/database/connection.go`

## 📝 Contribución

1. Fork el proyecto
2. Crear branch para la feature
3. Commit los cambios
4. Push al branch
5. Crear Pull Request

## 🚀 Deployment como Microservicio

Para información detallada sobre cómo desplegar este servicio como microservicio independiente, consulta la [Guía de Deployment](DEPLOYMENT.md).

### Configuraciones de ejemplo:

- **Desarrollo local**: Usar `docker-compose.yml` con variables de entorno locales
- **Producción**: Usar `examples/docker-compose.external.yml` para conectar a servicios remotos

## 📁 Archivos de Configuración

- `docker-compose.yml`: Configuración básica para desarrollo
- `examples/docker-compose.external.yml`: Ejemplo para conexión a servicios externos
- `DEPLOYMENT.md`: Guía detallada de deployment
- `examples/api_examples.md`: Ejemplos de uso de la API

## 📄 Licencia

Este proyecto está bajo la licencia MIT.
