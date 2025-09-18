package controllers

import (
	"holding-snapshots/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ValidationController struct {
	scrapingService *services.ScrapingService
}

// NewValidationController crea una nueva instancia del controlador
func NewValidationController() *ValidationController {
	return &ValidationController{
		scrapingService: services.NewScrapingService(),
	}
}

// ValidateHoldingRequest representa la estructura de la request de validación
type ValidateHoldingRequest struct {
	Code             string `json:"code" validate:"required"`
	TypeInvestmentID string `json:"typeInvestmentId,omitempty"`
}

// ValidateHoldingResponse representa la estructura de la response de validación
type ValidateHoldingResponse struct {
	IsValid bool `json:"isValid"`
}

// ValidateHolding valida si un holding es válido para ser agregado
// POST /api/validate
func (vc *ValidationController) ValidateHolding(c *fiber.Ctx) error {
	var req ValidateHoldingRequest

	// Parsear el body de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de request inválido",
		})
	}

	price, err := vc.scrapingService.ValidateHolding(req.TypeInvestmentID, req.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al validar el holding: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"holding": fiber.Map{
			"code":      req.Code,
			"lastPrice": price,
		},
		"isValid": price > 0,
	})
}

// HealthCheck endpoint simple para verificar el estado del servicio
// GET /api/health
func (vc *ValidationController) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "holding-snapshots",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
