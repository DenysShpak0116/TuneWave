package models

import "github.com/google/uuid"

type Author struct {
	BaseModel

	Name string `json:"name"`

	Songs []SongAuthor
}

type SongAuthor struct {
	BaseModel

	Role string `json:"role"`

	SongID   uuid.UUID `gorm:"constraint:OnDelete:CASCADE" json:"songId"`
	AuthorID uuid.UUID `gorm:"constraint:OnDelete:CASCADE" json:"authorId"`
	Author   Author
}
