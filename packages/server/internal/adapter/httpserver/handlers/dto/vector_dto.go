package dto

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type VectorDTO struct {
	ID uuid.UUID `json:"id"`

	Mark string `json:"mark"`

	CriterionID uuid.UUID `json:"criterionId"`
	Criterion   string    `json:"criterion"`

	CollectionSongID uuid.UUID `json:"collectionSongId"`
}

func (b *DTOBuilder) BuildVectorDTO(vector models.Vector) VectorDTO {
	return VectorDTO{
		ID:               vector.ID,
		Mark:             vector.Mark,
		CriterionID:      vector.CriterionID,
		Criterion:        vector.Criterion.Name,
		CollectionSongID: vector.CollectionSongID,
	}
}
