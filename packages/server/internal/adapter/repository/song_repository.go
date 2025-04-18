package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type SongRepository struct {
	*GenericRepository[models.Song]
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{
		GenericRepository: NewRepository[models.Song](db),
	}
}
