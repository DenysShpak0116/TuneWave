package models

import "github.com/google/uuid"

type Collection struct {
	BaseModel

	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`

	UserID uuid.UUID `json:"userId"`
	User   User      `gorm:"constraint:OnDelete:CASCADE"`

	CollectionSongs []CollectionSong `gorm:"foreignKey:CollectionID"`
	UserCollections []UserCollection `gorm:"foreignKey:CollectionID"`
}

type CollectionSong struct {
	BaseModel

	SongID uuid.UUID `json:"songId"`
	Song   Song

	CollectionID uuid.UUID  `json:"collectionId"`
	Collection   Collection `gorm:"constraint:OnDelete:CASCADE"`

	Vectors []Vector `gorm:"foreignKey:CollectionSongID"`
	Results []Result `gorm:"foreignKey:CollectionSongID"`
}
