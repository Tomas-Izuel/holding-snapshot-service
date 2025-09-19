package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Holding representa una tenencia específica de un activo en un grupo
type Holding struct {
	ID               string     `json:"id" gorm:"type:uuid;primary_key"`
	GroupID          string     `json:"groupId" gorm:"type:uuid;not null"`
	Group            Group      `json:"group" gorm:"foreignKey:GroupID"`
	Quantity         float64    `json:"quantity" gorm:"not null"`
	Earnings         float64    `json:"earnings" gorm:"not null;default:0"`
	RelativeEarnings float64    `json:"relativeEarnings" gorm:"not null;default:0"`
	AssetID          string     `json:"assetId" gorm:"type:uuid;not null"`
	Asset            Asset      `json:"asset" gorm:"foreignKey:AssetID"`
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
	return "Holding"
}

// CalculateEarnings calcula las ganancias basado en precio actual del activo vs snapshots anteriores
func (h *Holding) CalculateEarnings(currentPrice, previousPrice float64) {
	if previousPrice > 0 {
		earnings := (currentPrice - previousPrice) * h.Quantity
		h.Earnings = earnings

		relativeEarnings := ((currentPrice - previousPrice) / previousPrice) * 100
		h.RelativeEarnings = relativeEarnings
	}
}
