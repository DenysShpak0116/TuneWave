package dtos

import "github.com/google/uuid"

type VectorDTO struct {
	ID uuid.UUID `json:"id"`

	Mark string `json:"mark"`

	CriterionID uuid.UUID `json:"criterionId"`
	Criterion   string    `json:"criterion"`

	CollectionSongID uuid.UUID `json:"collectionSongId"`
}
