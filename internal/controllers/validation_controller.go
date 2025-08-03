package controllers

import (
	"holding-snapshots/internal/services"

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
	Name     string  `json:"name" validate:"required"`
	Code     string  `json:"code" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
	GroupID  string  `json:"groupId" validate:"required"`
}

// ValidateHoldingResponse representa la estructura de la response de validación
type ValidateHoldingResponse struct {
	Holding struct {
		Name     string  `json:"name"`
		Code     string  `json:"code"`
		Quantity float64 `json:"quantity"`
		GroupID  string  `json:"group_id"`
	} `json:"holding"`
	IsValid bool `json:"is_valid"`
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

	// Validaciones básicas
	if req.Name == "" || req.Code == "" || req.GroupID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Los campos name, code y groupId son requeridos",
		})
	}

	if req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "La cantidad debe ser mayor a 0",
		})
	}

	// Validar el holding usando el servicio de scraping
	holding, isValid, err := vc.scrapingService.ValidateHolding(
		req.Name, 
		req.Code, 
		req.GroupID, 
		req.Quantity,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error interno del servidor al validar el holding",
		})
	}

	// Construir response
	response := ValidateHoldingResponse{
		IsValid: isValid,
	}
	
	if holding != nil {
		response.Holding.Name = holding.Name
		response.Holding.Code = holding.Code
		response.Holding.Quantity = holding.Quantity
		response.Holding.GroupID = holding.GroupID
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// HealthCheck endpoint simple para verificar el estado del servicio
// GET /api/health
func (vc *ValidationController) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "holding-snapshots",
		"timestamp": fiber.Now(),
	})
}