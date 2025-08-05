package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TypeUser representa un tipo de usuario con permisos específicos
type TypeUser struct {
	ID          string       `json:"id" gorm:"type:uuid;primary_key"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Permissions []Permission `json:"permissions" gorm:"foreignKey:TypeUserID"`
	Users       []User       `json:"users" gorm:"foreignKey:TypeID"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (tu *TypeUser) BeforeCreate(tx *gorm.DB) error {
	if tu.ID == "" {
		tu.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (TypeUser) TableName() string {
	return "TypeUser"
}
