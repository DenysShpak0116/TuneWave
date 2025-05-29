package songservice

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

func (ss *SongService) AddToCollection(ctx context.Context, songUUID uuid.UUID, collectionUUID uuid.UUID) error {
	collectionSong := &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	}

	if err := ss.CollectionSongRepository.Add(ctx, collectionSong); err != nil {
		return err
	}

	return nil
}
