package database

import (
	"log"

	"holding-snapshots/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establece la conexión con la base de datos PostgreSQL
func Connect(databaseURL string) error {
	var err error
	
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(postgres.Open(databaseURL), config)
	if err != nil {
		return err
	}

	log.Println("✅ Conexión a la base de datos establecida")
	return nil
}

// AutoMigrate ejecuta las migraciones automáticas de GORM
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.TypeUser{},
		&models.Permission{},
		&models.User{},
		&models.TypeInvestment{},
		&models.Group{},
		&models.Holding{},
		&models.Snapshot{},
	)
	
	if err != nil {
		return err
	}

	log.Println("✅ Migraciones ejecutadas correctamente")
	return nil
}

// EnableUUIDExtension habilita la extensión uuid-ossp en PostgreSQL
func EnableUUIDExtension() error {
	err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return err
	}
	
	log.Println("✅ Extensión uuid-ossp habilitada")
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}