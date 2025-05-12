package models

import "github.com/google/uuid"

type Result struct {
	BaseModel

	SongRang int `json:"songRang"`

	UserID uuid.UUID `gorm:"type:uuid" json:"userId"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	CollectionSongID uuid.UUID      `gorm:"type:uuid" json:"collectionSongId"`
	CollectionSong   CollectionSong `gorm:"foreignKey:CollectionSongID;constraint:OnDelete:CASCADE;"`
}
