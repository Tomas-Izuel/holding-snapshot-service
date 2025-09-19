package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission representa un permiso específico del sistema
type Permission struct {
	ID          string   `json:"id" gorm:"type:uuid;primary_key"`
	Name        string   `json:"name" gorm:"uniqueIndex;not null"`
	Description string   `json:"description" gorm:"not null"`
	TypeUserID  string   `json:"typeUserId" gorm:"type:uuid;not null"`
	TypeUser    TypeUser `json:"typeUser" gorm:"foreignKey:TypeUserID"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (Permission) TableName() string {
	return "Permission"
}
