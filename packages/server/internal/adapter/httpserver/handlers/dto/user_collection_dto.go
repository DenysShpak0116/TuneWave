package dto

import (
	"time"

	"github.com/google/uuid"
)

type UsersCollectionDTO struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
}
