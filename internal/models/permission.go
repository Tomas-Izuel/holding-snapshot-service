package models

// Permission representa un permiso específico del sistema
type Permission struct {
	ID          string   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string   `json:"name" gorm:"uniqueIndex;not null"`
	Description string   `json:"description" gorm:"not null"`
	TypeUserID  string   `json:"type_user_id" gorm:"type:uuid;not null"`
	TypeUser    TypeUser `json:"type_user" gorm:"foreignKey:TypeUserID"`
}

// TableName especifica el nombre de la tabla
func (Permission) TableName() string {
	return "permissions"
}
