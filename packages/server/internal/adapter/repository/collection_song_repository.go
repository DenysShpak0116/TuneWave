package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type CollectionSongRepository struct {
	GenericRepository[models.CollectionSong]
}

func NewCollectionSongRepository(db *gorm.DB) *CollectionSongRepository {
	return &CollectionSongRepository{
		GenericRepository: GenericRepository[models.CollectionSong]{
			db: db,
		},
	}
}
