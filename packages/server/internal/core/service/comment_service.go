package service

import (
	"context"
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type CommentService struct {
	*GenericService[models.Comment]
}

func NewCommentService(repo port.Repository[models.Comment]) services.CommentService {
	return &CommentService{
		GenericService: NewGenericService(repo),
	}
}

func (cs *CommentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	comment, err := cs.Repository.NewQuery(ctx).Where("id = ?", id).
		Preload("User").Find()

	if err != nil {
		return nil, err
	}

	if len(comment) == 0 {
		return nil, fmt.Errorf("comment not found")
	}

	return &comment[0], nil
}
