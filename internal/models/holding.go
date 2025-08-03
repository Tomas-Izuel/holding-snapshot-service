package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Holding representa un activo de inversión específico
type Holding struct {
	ID               string     `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name             string     `json:"name" gorm:"not null"`              // Ej: "Apple", "BTC"
	Code             string     `json:"code" gorm:"not null"`              // Ej: "AAPL", "BTC"
	GroupID          string     `json:"group_id" gorm:"type:uuid;not null"`
	Group            Group      `json:"group" gorm:"foreignKey:GroupID"`
	Quantity         float64    `json:"quantity" gorm:"not null"`
	LastPrice        *float64   `json:"last_price,omitempty"`        // Puede ser null
	Earnings         *float64   `json:"earnings,omitempty"`          // Puede ser null
	RelativeEarnings *float64   `json:"relative_earnings,omitempty"` // Porcentaje de ganancia/pérdida
	Snapshots        []Snapshot `json:"snapshots" gorm:"foreignKey:HoldingID"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (h *Holding) BeforeCreate(tx *gorm.DB) error {
	if h.ID == "" {
		h.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (Holding) TableName() string {
	return "holdings"
}

// CalculateEarnings calcula las ganancias basado en precio actual vs snapshots anteriores
func (h *Holding) CalculateEarnings(currentPrice float64) {
	if h.LastPrice != nil && *h.LastPrice > 0 {
		earnings := (currentPrice - *h.LastPrice) * h.Quantity
		h.Earnings = &earnings
		
		relativeEarnings := ((currentPrice - *h.LastPrice) / *h.LastPrice) * 100
		h.RelativeEarnings = &relativeEarnings
	}
	h.LastPrice = &currentPrice
}