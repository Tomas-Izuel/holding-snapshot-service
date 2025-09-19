package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Group representa un grupo de inversión de un usuario
type Group struct {
	ID        string         `json:"id" gorm:"type:uuid;primary_key"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    string         `json:"userId" gorm:"type:uuid;not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	TypeID    string         `json:"typeId" gorm:"type:uuid;not null"`
	Type      TypeInvestment `json:"type" gorm:"foreignKey:TypeID"`
	Holdings  []Holding      `json:"holdings" gorm:"foreignKey:GroupID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

// TableName especifica el nombre de la tabla
func (Group) TableName() string {
	return "Group"
}
