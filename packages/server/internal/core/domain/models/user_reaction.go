package models

import "github.com/google/uuid"

type UserReaction struct {
	BaseModel

	Type   string    `json:"type"`
	UserID uuid.UUID `gorm:"type:uuid" json:"userId"`
	User   User      `gorm:"foreignKey:UserID"`

	SongID uuid.UUID `gorm:"type:uuid" json:"songId"`
	Song   Song      `gorm:"foreignKey:SongID"`
}
