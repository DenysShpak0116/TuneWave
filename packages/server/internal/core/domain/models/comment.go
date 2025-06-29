package models

import (
	"github.com/google/uuid"
)

type Comment struct {
	BaseModel

	Content string `json:"content"`

	UserID uuid.UUID `json:"userId"`
	User   User

	SongID uuid.UUID `json:"songId"`
}
