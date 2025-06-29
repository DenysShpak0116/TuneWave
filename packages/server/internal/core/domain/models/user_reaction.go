package models

import "github.com/google/uuid"

type UserReaction struct {
	BaseModel

	Type string `json:"type"`

	UserID uuid.UUID `json:"userId"`
	User   User

	SongID uuid.UUID `json:"songId"`
	Song   Song
}
