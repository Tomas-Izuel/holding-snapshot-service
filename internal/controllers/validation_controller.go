package controllers

import (
	"holding-snapshots/internal/models"
	"holding-snapshots/internal/services"
	"holding-snapshots/pkg/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ValidationController struct {
	scrapingService *services.ScrapingService
}

// NewValidationController crea una nueva instancia del controlador
func NewValidationController() *ValidationController {
	return &ValidationController{
		scrapingService: &services.ScrapingService{},
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

	// Validaciones básicas
	if req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El campo code es requerido",
		})
	}

	// Determinar si es validación con grupo existente o nuevo
	var holding *models.Holding
	var isValid bool
	var err error

	if req.GroupID != "" {
		// Validación con grupo existente
		quantity := req.Quantity
		if quantity == 0 {
			quantity = 1 // Valor por defecto para validación
		}

		log.Printf("🔍 Validación con grupo existente - GroupID: %s / Name: %s / Code: %s / Quantity: %f", req.GroupID, req.Name, req.Code, quantity)

		holding, isValid, err = vc.scrapingService.ValidateHolding(
			req.Name,
			req.Code,
			req.GroupID,
			quantity,
		)
	} else if req.TypeInvestmentID != "" && req.GroupName != "" {
		// Validación sin grupo existente (para creación de grupo)
		if req.TypeInvestmentID == "" || req.GroupName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Para validación sin grupo existente, typeInvestmentId y groupName son requeridos",
			})
		}

		// Obtener el tipo de inversión
		var typeInvestment models.TypeInvestment
		err = database.DB.First(&typeInvestment, "id = ?", req.TypeInvestmentID).Error
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tipo de inversión no encontrado",
			})
		}

		quantity := req.Quantity
		if quantity == 0 {
			quantity = 1 // Valor por defecto para validación
		}

		// Debug: Log del tipo de inversión encontrado
		log.Printf("🔍 Tipo de inversión encontrado: ID=%s, Name=%s, ScrapingURL=%s, Currency=%s",
			typeInvestment.ID, typeInvestment.Name, typeInvestment.ScrapingURL, typeInvestment.Currency)

		holding, isValid, err = vc.scrapingService.ValidateHoldingWithType(
			req.Name,
			req.Code,
			&typeInvestment,
			req.GroupName,
			quantity,
		)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Debe proporcionar groupId O (typeInvestmentId + groupName)",
		})
	}

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
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
