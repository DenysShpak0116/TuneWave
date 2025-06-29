package models

import "github.com/google/uuid"

type Vector struct {
	BaseModel

	Mark string `json:"mark"`

	CriterionID uuid.UUID `json:"criterionId"`
	Criterion   Criterion `gorm:"constraint:OnDelete:CASCADE"`

	CollectionSongID uuid.UUID `json:"collectionSongId"`
	CollectionSong   CollectionSong
}
