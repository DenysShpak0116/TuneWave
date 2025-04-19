package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CollectionDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
	CreatedAt   time.Time `json:"createdAt"`

	User UserDTO `json:"user"`

	CollectionSongs []SongPreviewDTO `json:"collectionSongs"`
}
