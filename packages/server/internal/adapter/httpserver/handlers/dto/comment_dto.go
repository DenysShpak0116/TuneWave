package dto

import (
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type CommentDTO struct {
	Author    UserDTO   `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
}

func (b *DTOBuilder) BuildCommentDTO(comment *models.Comment) CommentDTO {
	return CommentDTO{
		Author:    b.BuildUserDTO(&comment.User),
		CreatedAt: comment.CreatedAt,
		ID:        comment.ID,
		Content:   comment.Content,
	}
}
