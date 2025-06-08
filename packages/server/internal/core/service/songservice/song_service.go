package songservice

import (
	"bytes"
	"context"
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/google/uuid"
)

type SongService struct {
	*service.GenericService[models.Song]
	FileStorage              port.FileStorage
	AuthorRepository         port.Repository[models.Author]
	SongAuthorRepository     port.Repository[models.SongAuthor]
	TagsRepository           port.Repository[models.Tag]
	SongTagsRepository       port.Repository[models.SongTag]
	ReactionsRepository      port.Repository[models.UserReaction]
	CollectionSongRepository port.Repository[models.CollectionSong]
}

func NewSongService(
	songRepo port.Repository[models.Song],
	fileStorage port.FileStorage,
	authorRepo port.Repository[models.Author],
	songAuthorRepo port.Repository[models.SongAuthor],
	tagRepository port.Repository[models.Tag],
	songTagRepository port.Repository[models.SongTag],
	reactionsRepository port.Repository[models.UserReaction],
	collectionSongReporitory port.Repository[models.CollectionSong],
) services.SongService {
	return &SongService{
		GenericService:           service.NewGenericService(songRepo),
		FileStorage:              fileStorage,
		AuthorRepository:         authorRepo,
		SongAuthorRepository:     songAuthorRepo,
		TagsRepository:           tagRepository,
		SongTagsRepository:       songTagRepository,
		ReactionsRepository:      reactionsRepository,
		CollectionSongRepository: collectionSongReporitory,
	}
}

func (ss *SongService) GetSongs(ctx context.Context, search, sortBy, order string, page, limit int) ([]models.Song, error) {
	query := ss.Repository.NewQuery(ctx).
		Join("LEFT JOIN song_authors ON song_authors.song_id = songs.id").
		Join("LEFT JOIN authors ON authors.id = song_authors.author_id").
		Join("LEFT JOIN song_tags ON song_tags.song_id = songs.id").
		Join("LEFT JOIN tags ON tags.id = song_tags.tag_id").
		Where("songs.title ILIKE ? OR authors.name ILIKE ? OR LOWER(songs.genre) = LOWER(?)",
			"%"+search+"%", "%"+search+"%", search).
		Order(fmt.Sprintf("songs.%s %s", sortBy, order)).
		Group("songs.id").
		Skip((page - 1) * limit).
		Take(limit).
		Preload("User").
		Preload("Authors").
		Preload("Authors.Author").
		Preload("SongTags").
		Preload("SongTags.Tag").
		Preload("Comments").
		Preload("Comments.User")

	songs, err := query.Find()
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return []models.Song{}, nil
	}

	return songs, nil
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil
}

func (ss *SongService) ReactionsCount(ctx context.Context, id uuid.UUID, reactionType string) (int64, error) {
	reactionAmount, err := ss.ReactionsRepository.NewQuery(ctx).
		Where("song_id = ? AND type = ?", id, reactionType).
		Count()
	if err != nil {
		return 0, err
	}
	return reactionAmount, nil
}
