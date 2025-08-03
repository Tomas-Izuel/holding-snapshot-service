package models

// TypeUser representa un tipo de usuario con permisos específicos
type TypeUser struct {
	ID          string       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Permissions []Permission `json:"permissions" gorm:"foreignKey:TypeUserID"`
	Users       []User       `json:"users" gorm:"foreignKey:TypeID"`
}

// TableName especifica el nombre de la tabla
func (TypeUser) TableName() string {
	return "type_users"
}
