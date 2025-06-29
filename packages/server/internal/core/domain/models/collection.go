package models

import "github.com/google/uuid"

type Collection struct {
	BaseModel

	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`

	UserID uuid.UUID `json:"userId"`
	User   User

	CollectionSongs []CollectionSong `gorm:"constraint:OnDelete:CASCADE"`
	UserCollections []UserCollection `gorm:"constraint:OnDelete:CASCADE"`
}

type CollectionSong struct {
	BaseModel

	SongID uuid.UUID `json:"songId"`
	Song   Song

	CollectionID uuid.UUID `json:"collectionId"`
	Collection   Collection

	Vectors []Vector `gorm:"constraint:OnDelete:CASCADE"`
	Results []Result `gorm:"constraint:OnDelete:CASCADE"`
}
