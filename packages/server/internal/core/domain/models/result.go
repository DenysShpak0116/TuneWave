package models

import "github.com/google/uuid"

type Result struct {
	BaseModel

	SongRang int `json:"songRang"`

	UserID uuid.UUID `json:"userId"`
	User   User      `gorm:"constraint:OnDelete:CASCADE;"`

	CollectionSongID uuid.UUID      `json:"collectionSongId"`
	CollectionSong   CollectionSong `gorm:"constraint:OnDelete:CASCADE;"`
}
