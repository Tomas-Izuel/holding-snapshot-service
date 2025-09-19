package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Asset representa un activo de inversión
type Asset struct {
	ID        string         `json:"id" gorm:"type:uuid;primary_key"`
	Name      string         `json:"name" gorm:"not null"`
	Code      string         `json:"code" gorm:"not null;uniqueIndex"`
	LastPrice float64        `json:"lastPrice" gorm:"not null;column:lastPrice"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:createdAt"`
	TypeID    string         `json:"typeId" gorm:"type:uuid;not null;column:typeId"`
	IsValid   bool           `json:"isValid" gorm:"default:true;column:is_valid"`
	Type      TypeInvestment `json:"type" gorm:"foreignKey:TypeID;references:ID"`
	Holdings  []Holding      `json:"holdings" gorm:"foreignKey:AssetID"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (a *Asset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (Asset) TableName() string {
	return "Asset"
}
