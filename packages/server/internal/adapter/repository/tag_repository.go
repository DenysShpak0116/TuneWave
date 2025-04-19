package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type TagRepository struct {
	GenericRepository[models.Tag]
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		GenericRepository: GenericRepository[models.Tag]{
			db: db,
		},
	}
}
