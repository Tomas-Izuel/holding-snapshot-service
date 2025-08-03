package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	RedisURL          string
	Port              string
	Environment       string
	APIKey            string
	ScrapingCronSchedule string
}

var AppConfig *Config

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	config := &Config{
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		RedisURL:            getEnv("REDIS_URL", "redis://localhost:6379"),
		Port:                getEnv("PORT", "8080"),
		Environment:         getEnv("ENV", "development"),
		APIKey:              getEnv("SNAPSHOT_SERVICE_API_KEY", ""),
		ScrapingCronSchedule: getEnv("SCRAPING_CRON_SCHEDULE", "0 1 * * 0"), // Domingos 1:00 AM
	}

	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL es requerido")
	}

	if config.APIKey == "" {
		log.Fatal("SNAPSHOT_SERVICE_API_KEY es requerido")
	}

	AppConfig = config
	return config
}

// getEnv obtiene una variable de entorno con un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}