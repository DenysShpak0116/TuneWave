package dto

import (
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type UserCollectionDTO struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
}

func (b *DTOBuilder) BuildUserCollectionDTO(collection *models.Collection) UserCollectionDTO {
	return UserCollectionDTO{
		ID:          collection.ID,
		CreatedAt:   collection.CreatedAt,
		Title:       collection.Title,
		Description: collection.Description,
		CoverURL:    collection.CoverURL,
	}
}
