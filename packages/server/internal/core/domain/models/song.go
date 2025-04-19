package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	BaseModel

	Title      string        `json:"title"`
	Genre      string        `json:"genre"`
	SongURL    string        `json:"songUrl"`
	CoverURL   string        `json:"coverUrl"`
	Duration   time.Duration `json:"duration"`
	Listenings int64         `json:"listenings"`

	UserID uuid.UUID `gorm:"type:uuid" json:"userId"`
	User   User      `gorm:"foreignKey:UserID"`

	SongTags        []SongTag    `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE"`
	Comments        []Comment    `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE"`
	Authors         []SongAuthor `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE"`
	CollectionSongs []CollectionSong
	Reactions       []UserReaction `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE"`
}
