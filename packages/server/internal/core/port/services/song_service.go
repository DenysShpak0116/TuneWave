package services

import (
	"context"
	"mime/multipart"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type SaveSongParams struct {
	UserID      uuid.UUID
	Title       string
	Genre       string
	Artists     []string
	Tags        []string
	Song        multipart.File
	SongHeader  *multipart.FileHeader
	Cover       multipart.File
	CoverHeader *multipart.FileHeader
}

type SongService interface {
	Service[models.Song]
	GetSongs(
		ctx context.Context,
		search string,
		sortBy string,
		order string,
		page int,
		limit int,
	) ([]dtos.SongDTO, error)
	SaveSong(ctx context.Context, songParams SaveSongParams) (*models.Song, error)
	ReactionsCount(ctx context.Context, id uuid.UUID, reactionType string) (int64, error)
	GetFullDTOByID(ctx context.Context, id uuid.UUID) (*dtos.SongExtendedDTO, error)
	SetReaction(ctx context.Context, songID uuid.UUID, userID uuid.UUID, reactionType string) (int64, int64, error)
	AddToCollection(ctx context.Context, songUUID, collectionUUID uuid.UUID) error
}
