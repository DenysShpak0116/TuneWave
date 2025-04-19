package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type SongAuthorRepository struct {
	GenericRepository[models.SongAuthor]
}

func NewSongAuthorRepository(db *gorm.DB) *SongAuthorRepository {
	return &SongAuthorRepository{
		GenericRepository: GenericRepository[models.SongAuthor]{
			db: db,
		},
	}
}
