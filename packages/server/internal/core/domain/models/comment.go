package models

import (
	"github.com/google/uuid"
)

type Comment struct {
	BaseModel

	Content string `json:"content"`

	UserID uuid.UUID `gorm:"type:uuid" json:"userId"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	SongID uuid.UUID `gorm:"type:uuid" json:"songId"`
	Song   Song      `gorm:"foreignKey:SongID;constraint:OnDelete:CASCADE"`
}
