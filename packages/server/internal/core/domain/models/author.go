package models

import "github.com/google/uuid"

type Author struct {
	BaseModel
	Name  string       `json:"name"`
	Songs []SongAuthor `gorm:"constraint:OnDelete:CASCADE"`
}

type SongAuthor struct {
	BaseModel

	Role string `json:"role"`

	SongID   uuid.UUID `json:"songId"`
	AuthorID uuid.UUID `json:"authorId"`
	Author   Author
}
