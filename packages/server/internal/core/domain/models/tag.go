package models

import "github.com/google/uuid"

type Tag struct {
	BaseModel

	Name string `json:"name"`

	SongTags []SongTag
}

type SongTag struct {
	BaseModel

	TagID uuid.UUID `gorm:"type:uuid" json:"tagId"`
	Tag   Tag       `gorm:"foreignKey:TagID"`

	SongID uuid.UUID `gorm:"type:uuid" json:"songId"`
	Song   Song      `gorm:"foreignKey:SongID"`
}
