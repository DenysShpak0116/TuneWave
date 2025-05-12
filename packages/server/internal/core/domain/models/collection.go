package models

import "github.com/google/uuid"

type Collection struct {
	BaseModel

	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`

	UserID uuid.UUID `gorm:"type:uuid" json:"userId"`
	User   User      `gorm:"foreignKey:UserID"`

	CollectionSongs []CollectionSong `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE"`
	UserCollections []UserCollection `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE"`
}

type CollectionSong struct {
	BaseModel

	SongID uuid.UUID `gorm:"type:uuid" json:"songId"`
	Song   Song      `gorm:"foreignKey:SongID"`

	CollectionID uuid.UUID  `gorm:"type:uuid" json:"collectionId"`
	Collection   Collection `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE"`

	Vectors []Vector `gorm:"foreignKey:CollectionSongID;constraint:OnDelete:CASCADE"`
	Results []Result `gorm:"foreignKey:CollectionSongID;constraint:OnDelete:CASCADE"`
}
