package repository

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"gorm.io/gorm"
)

type CommentRepository struct {
	GenericRepository[models.Comment]
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		GenericRepository: GenericRepository[models.Comment]{
			db: db,
		},
	}
}
