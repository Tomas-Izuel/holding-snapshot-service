package utils

import (
	"strings"

	"github.com/google/uuid"
)

// IsValidUUID valida si un string es un UUID válido
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// IsValidCode valida si un código de activo es válido
func IsValidCode(code string) bool {
	if len(code) < 1 || len(code) > 10 {
		return false
	}
	
	// Solo permitir letras, números y algunos caracteres especiales
	return strings.ContainsAny(code, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.")
}

// IsValidQuantity valida si una cantidad es válida
func IsValidQuantity(quantity float64) bool {
	return quantity > 0 && quantity <= 1000000 // Límite razonable
}

// SanitizeString limpia y valida strings de entrada
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}