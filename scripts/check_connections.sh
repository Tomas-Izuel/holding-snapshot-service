#!/bin/bash

# Script para verificar conexiones a servicios externos antes del deployment

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔍 Verificando conexiones a servicios externos...${NC}"
echo "=================================================="

# Verificar que las variables de entorno estén configuradas
echo -e "\n${YELLOW}1. Verificando variables de entorno...${NC}"

if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}❌ DATABASE_URL no está configurada${NC}"
    exit 1
else
    echo -e "${GREEN}✅ DATABASE_URL configurada${NC}"
fi

if [ -z "$REDIS_URL" ]; then
    echo -e "${RED}❌ REDIS_URL no está configurada${NC}"
    exit 1
else
    echo -e "${GREEN}✅ REDIS_URL configurada${NC}"
fi

if [ -z "$SNAPSHOT_SERVICE_API_KEY" ]; then
    echo -e "${RED}❌ SNAPSHOT_SERVICE_API_KEY no está configurada${NC}"
    exit 1
else
    echo -e "${GREEN}✅ SNAPSHOT_SERVICE_API_KEY configurada${NC}"
fi

# Extraer host y puerto de DATABASE_URL
DB_HOST=$(echo $DATABASE_URL | sed -n 's/.*@\([^:]*\):.*/\1/p')
DB_PORT=$(echo $DATABASE_URL | sed -n 's/.*:\([0-9]*\)\/.*/\1/p')

# Extraer host y puerto de REDIS_URL
REDIS_HOST=$(echo $REDIS_URL | sed -n 's/.*:\/\/\([^:]*\):.*/\1/p')
REDIS_PORT=$(echo $REDIS_URL | sed -n 's/.*:\([0-9]*\)$/\1/p')

echo -e "\n${YELLOW}2. Verificando conectividad de red...${NC}"

# Verificar conectividad a PostgreSQL
echo -n "Verificando PostgreSQL ($DB_HOST:$DB_PORT)... "
if timeout 5 bash -c "</dev/tcp/$DB_HOST/$DB_PORT" 2>/dev/null; then
    echo -e "${GREEN}✅ Accesible${NC}"
else
    echo -e "${RED}❌ No accesible${NC}"
    exit 1
fi

# Verificar conectividad a Redis
echo -n "Verificando Redis ($REDIS_HOST:$REDIS_PORT)... "
if timeout 5 bash -c "</dev/tcp/$REDIS_HOST/$REDIS_PORT" 2>/dev/null; then
    echo -e "${GREEN}✅ Accesible${NC}"
else
    echo -e "${RED}❌ No accesible${NC}"
    exit 1
fi

echo -e "\n${YELLOW}3. Verificando acceso a bases de datos...${NC}"

# Verificar PostgreSQL con psql si está disponible
if command -v psql &> /dev/null; then
    echo -n "Verificando autenticación PostgreSQL... "
    if psql "$DATABASE_URL" -c "SELECT 1;" &> /dev/null; then
        echo -e "${GREEN}✅ Conexión exitosa${NC}"
    else
        echo -e "${RED}❌ Error de autenticación${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  psql no disponible, saltando test de autenticación PostgreSQL${NC}"
fi

# Verificar Redis con redis-cli si está disponible
if command -v redis-cli &> /dev/null; then
    echo -n "Verificando conexión Redis... "
    if redis-cli -u "$REDIS_URL" ping &> /dev/null; then
        echo -e "${GREEN}✅ Conexión exitosa${NC}"
    else
        echo -e "${RED}❌ Error de conexión${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  redis-cli no disponible, saltando test de conexión Redis${NC}"
fi

echo -e "\n${GREEN}🎉 Todas las verificaciones pasaron. El servicio está listo para deployment.${NC}"
echo "=================================================="