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

	UserID uuid.UUID `json:"userId"`
	User   User

	SongTags        []SongTag        `gorm:"constraint:OnDelete:CASCADE"`
	Comments        []Comment        `gorm:"constraint:OnDelete:CASCADE"`
	Authors         []SongAuthor     `gorm:"constraint:OnDelete:CASCADE"`
	CollectionSongs []CollectionSong `gorm:"constraint:OnDelete:CASCADE"`
	Reactions       []UserReaction   `gorm:"constraint:OnDelete:CASCADE"`
}
