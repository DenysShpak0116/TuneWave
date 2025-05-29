package models

import "github.com/google/uuid"

type Tag struct {
	BaseModel

	Name string `json:"name"`

	SongTags []SongTag
}

type SongTag struct {
	BaseModel

	TagID uuid.UUID `json:"tagId"`
	Tag   Tag

	SongID uuid.UUID `json:"songId"`
	Song   Song
}
