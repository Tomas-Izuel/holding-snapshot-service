package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Snapshot representa un snapshot del precio de un holding en un momento específico
type Snapshot struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	Price     float64   `json:"price" gorm:"not null"` // Precio del holding en el momento del snapshot
	HoldingID string    `json:"holding_id" gorm:"type:uuid;not null"`
	Holding   Holding   `json:"holding" gorm:"foreignKey:HoldingID"`
	Quantity  float64   `json:"quantity" gorm:"not null"` // Cantidad de holdings al momento del snapshot
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (s *Snapshot) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (Snapshot) TableName() string {
	return "Snapshot"
}
