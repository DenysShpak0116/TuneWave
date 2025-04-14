package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGORMDB(connectionString string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate()
}
