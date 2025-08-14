package models

import "github.com/google/uuid"

type Result struct {
	BaseModel

	SongRank int `json:"songRank"`

	UserID uuid.UUID `json:"userId"`
	User   User

	CollectionSongID uuid.UUID `json:"collectionSongId"`
	CollectionSong   CollectionSong
}
