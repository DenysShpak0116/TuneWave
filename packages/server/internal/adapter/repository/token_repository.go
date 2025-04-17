package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type TokenRepository struct {
	GenericRepository[models.Token]
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{
		GenericRepository: GenericRepository[models.Token]{db: db},
	}
}
