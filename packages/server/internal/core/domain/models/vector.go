package models

import "github.com/google/uuid"

type Vector struct {
	BaseModel

	Mark string `json:"mark"`

	CriterionID uuid.UUID `gorm:"type:uuid" json:"criterionId"`
	Criterion   Criterion `gorm:"foreignKey:CriterionID"`

	CollectionSongID uuid.UUID      `gorm:"type:uuid" json:"collectionSongId"`
	CollectionSong   CollectionSong `gorm:"foreignKey:CollectionSongID"`
}
