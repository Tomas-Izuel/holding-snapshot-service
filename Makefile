# Makefile para el proyecto Holding Snapshots

.PHONY: help build run test clean docker-build docker-run docker-down deps lint

# Variables
APP_NAME=holding-snapshots
DOCKER_IMAGE=$(APP_NAME):latest
GO_VERSION=1.21

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  build         - Compilar la aplicación"
	@echo "  run           - Ejecutar la aplicación"
	@echo "  test          - Ejecutar tests"
	@echo "  clean         - Limpiar archivos compilados"
	@echo "  deps          - Instalar/actualizar dependencias"
	@echo "  lint          - Ejecutar linters"
	@echo "  docker-build  - Construir imagen Docker"
	@echo "  docker-run    - Ejecutar con Docker Compose"
	@echo "  docker-down   - Parar contenedores Docker"

# Compilar la aplicación
build:
	@echo "🔨 Compilando $(APP_NAME)..."
	go build -o bin/$(APP_NAME) cmd/server/main.go
	@echo "✅ Compilación completada"

# Ejecutar la aplicación
run:
	@echo "🚀 Ejecutando $(APP_NAME)..."
	go run cmd/server/main.go

# Ejecutar tests
test:
	@echo "🧪 Ejecutando tests..."
	go test -v ./...

# Limpiar archivos compilados
clean:
	@echo "🧹 Limpiando archivos compilados..."
	rm -rf bin/
	go clean

# Instalar/actualizar dependencias
deps:
	@echo "📦 Instalando dependencias..."
	go mod download
	go mod tidy

# Ejecutar linters
lint:
	@echo "🔍 Ejecutando linters..."
	go fmt ./...
	go vet ./...

# Construir imagen Docker
docker-build:
	@echo "🐳 Construyendo imagen Docker..."
	docker build -t $(DOCKER_IMAGE) .

# Ejecutar con Docker Compose
docker-run:
	@echo "🐳 Iniciando servicio holding-snapshots..."
	docker-compose up --build -d

# Parar contenedores Docker
docker-down:
	@echo "🐳 Parando servicio holding-snapshots..."
	docker-compose down

# Inicializar el proyecto
init: deps
	@echo "🎯 Inicializando proyecto..."
	@if [ ! -f .env ]; then \
		echo "📋 Creando archivo .env..."; \
		echo "DATABASE_URL=postgresql://usuario:password@localhost:5432/dbname" > .env; \
		echo "REDIS_URL=redis://localhost:6379" >> .env; \
		echo "SNAPSHOT_SERVICE_API_KEY=tu_api_key_aqui" >> .env; \
		echo "PORT=8080" >> .env; \
		echo "ENV=development" >> .env; \
		echo "SCRAPING_CRON_SCHEDULE=0 1 * * 0" >> .env; \
		echo "⚠️  Configura las variables en .env para conectar con el servicio principal"; \
	fi
	@echo "✅ Proyecto inicializado"

# Ejecutar migraciones
migrate:
	@echo "📊 Ejecutando migraciones..."
	go run cmd/server/main.go --migrate-only

# Ver logs de Docker Compose
logs:
	docker-compose logs -f

# Entrar al contenedor de la aplicación
shell:
	docker-compose exec holding-snapshots sh

# Verificar conexiones antes del deployment
check-connections:
	@echo "🔍 Verificando conexiones a servicios externos..."
	@chmod +x scripts/check_connections.sh
	@./scripts/check_connections.sh

# Deployment con verificaciones
deploy: check-connections docker-build docker-run
	@echo "🚀 Deployment completado"
	@echo "Verificando salud del servicio..."
	@sleep 5
	@curl -f http://localhost:8080/api/health || echo "⚠️  El servicio puede necesitar más tiempo para iniciar"