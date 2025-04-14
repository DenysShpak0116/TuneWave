package repository

import (
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGORMDB(connectionString string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	defer fmt.Println("Migration completed")
	return db.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.Message{},
		&models.Song{},
		&models.Author{},
		&models.SongAuthor{},
		&models.Collection{},
		&models.CollectionSong{},
		&models.Comment{},
		&models.UserReaction{},
		&models.Tag{},
		&models.SongTag{},
	)
}
