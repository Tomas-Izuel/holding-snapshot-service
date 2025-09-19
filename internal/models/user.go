package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User representa un usuario del sistema
type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Password  string    `json:"-" gorm:"not null"` // No incluir en JSON por seguridad
	TypeID    string    `json:"typeId" gorm:"type:uuid;not null"`
	Type      TypeUser  `json:"type" gorm:"foreignKey:TypeID"`
	Groups    []Group   `json:"groups" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (User) TableName() string {
	return "User"
}
