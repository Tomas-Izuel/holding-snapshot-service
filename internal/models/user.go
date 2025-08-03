package models

import (
	"time"
)

// User representa un usuario del sistema
type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Password  string    `json:"-" gorm:"not null"` // No incluir en JSON por seguridad
	TypeID    string    `json:"type_id" gorm:"type:uuid;not null"`
	Type      TypeUser  `json:"type" gorm:"foreignKey:TypeID"`
	Groups    []Group   `json:"groups" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (User) TableName() string {
	return "users"
}
