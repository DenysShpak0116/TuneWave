package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGORMDB(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.StoragePath), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Warn,
				Colorful:      true,
			},
		),
	})
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
		&models.Token{},
		&models.Criterion{},
		&models.Vector{},
		&models.Result{},
		&models.UserCollection{},
		&models.UserFollower{},
	)
}
