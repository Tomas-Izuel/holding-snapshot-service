package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TypeInvestment representa un tipo de inversión con su URL de scraping
type TypeInvestment struct {
	ID          string  `json:"id" gorm:"type:uuid;primary_key"`
	Name        string  `json:"name" gorm:"not null"` // Ej: "Cedears", "Criptomonedas", "Acciones"
	ScrapingURL string  `json:"scrappingUrl" gorm:"column:scrappingUrl;not null"`
	Currency    string  `json:"currency" gorm:"not null"` // Ej: "USD", "ARS"
	Groups      []Group `json:"groups" gorm:"foreignKey:TypeID"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (ti *TypeInvestment) BeforeCreate(tx *gorm.DB) error {
	if ti.ID == "" {
		ti.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (TypeInvestment) TableName() string {
	return "TypeInvestment"
}
