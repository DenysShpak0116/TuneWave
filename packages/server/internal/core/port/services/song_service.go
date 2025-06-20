package services

import (
	"context"
	"mime/multipart"

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

type UpdateSongParams struct {
	SongID      uuid.UUID
	Title       string
	Genre       string
	Artists     []string
	Tags        []string
	Song        multipart.File
	SongHeader  *multipart.FileHeader
	Cover       multipart.File
	CoverHeader *multipart.FileHeader
}

type SearchSongsParams struct {
	Search string
	SortBy string
	Order  string
	Page   int
	Limit  int
}

type SongService interface {
	Service[models.Song]
	GetSongs(ctx context.Context, params SearchSongsParams, preloads ...string) ([]models.Song, error)
	SaveSong(ctx context.Context, songParams SaveSongParams) (*models.Song, error)
	UpdateSong(ctx context.Context, songParams UpdateSongParams) error
	ReactionsCount(ctx context.Context, id uuid.UUID, reactionType string) (int64, error)
	SetReaction(ctx context.Context, songID uuid.UUID, userID uuid.UUID, reactionType string) (int64, int64, error)
	AddToCollection(ctx context.Context, songUUID, collectionUUID uuid.UUID) error
	IsReactedByUser(ctx context.Context, songID uuid.UUID, userID uuid.UUID) (string, error)
	GetGenres(context.Context) []string
	GetGenresMostPopularSong(ctx context.Context, genre string) (*models.Song, error)
}
