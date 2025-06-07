package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CollectionDTO struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	CoverURL string    `json:"coverUrl"`

	User UserDTO `json:"user"`
}

type UsersCollectionDTO struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
}
