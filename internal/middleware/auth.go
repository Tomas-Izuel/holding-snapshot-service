package middleware

import (
	"holding-snapshots/internal/config"

	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth middleware para validar la API Key
func APIKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("Authorization")
		
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API Key requerida en el header Authorization",
			})
		}

		// Validar API Key
		if apiKey != config.AppConfig.APIKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API Key inválida",
			})
		}

		return c.Next()
	}
}

// CORS middleware personalizado
func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}