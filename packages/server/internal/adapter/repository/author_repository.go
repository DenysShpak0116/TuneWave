package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	GenericRepository[models.Author]
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{
		GenericRepository: GenericRepository[models.Author]{
			db: db,
		},
	}
}
