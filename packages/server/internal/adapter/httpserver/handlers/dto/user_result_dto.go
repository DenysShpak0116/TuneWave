package dto

import "github.com/google/uuid"

type UserResultsDTO struct {
	CollectionSongID uuid.UUID `json:"collectionSongId"`

	SongID   uuid.UUID `json:"songId"`
	SongName string    `json:"songName"`

	UserID   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
	SongRang int       `json:"songRang"`
}
