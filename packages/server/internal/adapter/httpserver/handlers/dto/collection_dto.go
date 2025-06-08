package dto

import (
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type CollectionDTO struct {
	CreatedAt   time.Time `json:"createdAt"`
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	CoverURL    string    `json:"coverUrl"`
	Description string    `json:"description"`

	User *UserDTO `json:"user"`

	CollectionSongs []SongDTO `json:"collectionSongs"`
}

func NewCollectionDTO(collection *models.Collection) *CollectionDTO {
	collectionSongs := make([]SongDTO, 0)
	for _, collectionSong := range collection.CollectionSongs {
		collectionSongs = append(collectionSongs, *NewSongDTO(&collectionSong.Song))
	}

	return &CollectionDTO{
		ID:              collection.ID,
		Title:           collection.Title,
		CoverURL:        collection.CoverURL,
		CreatedAt:       collection.CreatedAt,
		Description:     collection.Description,
		User:            NewUserDTO(&collection.User),
		CollectionSongs: collectionSongs,
	}
}
