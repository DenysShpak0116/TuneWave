package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	GenericRepository[models.Collection]
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{
		GenericRepository: GenericRepository[models.Collection]{
			db: db,
		},
	}
}
