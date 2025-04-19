package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type SongTagRepository struct {
	GenericRepository[models.SongTag]
}

func NewSongTagRepository(db *gorm.DB) *SongTagRepository {
	return &SongTagRepository{
		GenericRepository: GenericRepository[models.SongTag]{
			db: db,
		},
	}
}
