package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	BaseModel

	Title      string    `json:"title"`
	ArtistName string    `json:"artistName"`
	UploadDate time.Time `json:"uploadDate"`
	Genre      string    `json:"genre"`
	FileURL    string    `json:"fileUrl"`
	CoverURL   string    `json:"coverUrl"`
	Duration   float64   `json:"duration"`
	Listenings int       `json:"listenings"`

	UserID uuid.UUID `gorm:"type:uuid" json:"userId"`
	User   User      `gorm:"foreignKey:UserID"`

	SongTags        []SongTag
	Comments        []Comment
	Authors         []SongAuthor
	CollectionSongs []CollectionSong
}
