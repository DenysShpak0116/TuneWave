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

	User UserDTO `json:"user"`

	CollectionSongs []SongDTO `json:"collectionSongs,omitempty"`
}

func (b *DTOBuilder) BuildCollectionDTO(collection *models.Collection) CollectionDTO {
	collectionSongs := []SongDTO{}
	for _, collectionSong := range collection.CollectionSongs {
		collectionSongs = append(collectionSongs, *b.BuildSongDTO(&collectionSong.Song))
	}

	if len(collectionSongs) == 0 {
		collectionSongs = nil
	}

	return CollectionDTO{
		ID:              collection.ID,
		Title:           collection.Title,
		CoverURL:        collection.CoverURL,
		CreatedAt:       collection.CreatedAt,
		Description:     collection.Description,
		User:            b.BuildUserDTO(&collection.User),
		CollectionSongs: collectionSongs,
	}
}
