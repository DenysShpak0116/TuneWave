package songservice

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	dtoMapper "github.com/dranikpg/dto-mapper"
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

func (ss *SongService) GetSongs(ctx context.Context, search, sortBy, order string, page, limit int) ([]dtos.SongDTO, error) {
	songs, err := ss.Repository.NewQuery(ctx).
		Where("title LIKE ?", "%"+search+"%").
		Order(fmt.Sprintf("%s %s", sortBy, order)).
		Skip((page - 1) * limit).
		Take(limit).
		Preload("User").
		Find()
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return nil, fmt.Errorf("no songs found for search term: %s", search)
	}

	songDTOs := make([]dtos.SongDTO, len(songs))
	for i, song := range songs {
		likes, err := ss.ReactionsRepository.NewQuery(ctx).
			Where("song_id = ? AND type = ?", song.ID, "like").
			Count()
		if err != nil {
			return nil, fmt.Errorf("get likes: %w", err)
		}

		dislikes, err := ss.ReactionsRepository.NewQuery(ctx).
			Where("song_id = ? AND type = ?", song.ID, "dislike").
			Count()
		if err != nil {
			return nil, fmt.Errorf("get dislikes: %w", err)
		}

		songDTOs[i] = dtos.SongDTO{
			ID:         song.ID,
			Duration:   formatDuration(song.Duration),
			Title:      song.Title,
			SongURL:    song.SongURL,
			CoverURL:   song.CoverURL,
			Listenings: song.Listenings,
			Likes:      likes,
			Dislikes:   dislikes,
			User: dtos.UserDTO{
				ID:             song.User.ID,
				Username:       song.User.Username,
				ProfilePicture: song.User.ProfilePicture,
			},
		}
	}
	return songDTOs, nil
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil
}

func (ss *SongService) GetByID(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	songs, err := ss.Repository.NewQuery(ctx).
		Where("id = ?", id).
		Preload("Authors").
		Preload("Authors.Author").
		Preload("SongTags").
		Preload("SongTags.Tag").
		Preload("SongTags").
		Preload("Comments").
		Preload("Comments.User").
		Preload("User").
		Preload("Reactions").
		Find()
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return nil, fmt.Errorf("song not found")
	}
	return &songs[0], nil
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

func (ss *SongService) GetFullDTOByID(ctx context.Context, id uuid.UUID) (*dtos.SongExtendedDTO, error) {
	song, err := ss.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var songUserDTO dtos.UserDTO
	if err := dtoMapper.Map(&songUserDTO, song.User); err != nil {
		return nil, fmt.Errorf("map user: %w", err)
	}

	authorsDTO := make([]dtos.AuthorDTO, 0)
	for _, sa := range song.Authors {
		authorsDTO = append(authorsDTO, dtos.AuthorDTO{
			ID:   sa.ID,
			Name: sa.Author.Name,
			Role: sa.Role,
		})
	}

	tagsDTO := make([]dtos.TagDTO, 0)
	for _, st := range song.SongTags {
		tagsDTO = append(tagsDTO, dtos.TagDTO{
			Name: st.Tag.Name,
		})
	}

	commentsDTO := make([]dtos.CommentDTO, 0)
	for _, c := range song.Comments {
		commentsDTO = append(commentsDTO, dtos.CommentDTO{
			ID: c.ID,
			Author: dtos.UserDTO{
				ID:             c.User.ID,
				Username:       c.User.Username,
				ProfilePicture: c.User.ProfilePicture,
			},
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		})
	}

	likes, err := ss.ReactionsCount(ctx, id, "like")
	if err != nil {
		return nil, fmt.Errorf("get likes: %w", err)
	}
	dislikes, err := ss.ReactionsCount(ctx, id, "dislike")
	if err != nil {
		return nil, fmt.Errorf("get dislikes: %w", err)
	}

	dto := &dtos.SongExtendedDTO{
		ID:         song.ID,
		Title:      song.Title,
		Genre:      song.Genre,
		SongURL:    song.SongURL,
		CoverURL:   song.CoverURL,
		CreatedAt:  song.CreatedAt,
		Duration:   formatDuration(time.Duration(song.Duration)),
		Listenings: song.Listenings,
		User:       songUserDTO,
		Authors:    authorsDTO,
		SongTags:   tagsDTO,
		Comments:   commentsDTO,
		Likes:      likes,
		Dislikes:   dislikes,
	}

	return dto, nil
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
