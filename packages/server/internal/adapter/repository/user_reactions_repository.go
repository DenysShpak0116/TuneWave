package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type UserReactionRepository struct {
	GenericRepository[models.UserReaction]
}

func NewUserReactionRepository(db *gorm.DB) *UserReactionRepository {
	return &UserReactionRepository{
		GenericRepository: GenericRepository[models.UserReaction]{

			db: db,
		},
	}
}
