package dtos

import "github.com/google/uuid"

type SongPreviewDTO struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Duration   string    `json:"duration"`
	CoverURL   string    `json:"coverUrl"`
	Listenings int64     `json:"listenings"`
	UserID     uuid.UUID `json:"userId"`
	UserName   string    `json:"userName"`
}
