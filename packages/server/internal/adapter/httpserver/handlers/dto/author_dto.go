package dto

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type AuthorDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role string    `json:"role"`
}

func (b *DTOBuilder) BuildAuthorDTO(songAuthor *models.SongAuthor) *AuthorDTO {
	return &AuthorDTO{
		ID:   songAuthor.AuthorID,
		Name: songAuthor.Author.Name,
		Role: songAuthor.Role,
	}
}
