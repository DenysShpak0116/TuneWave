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

func (ss *SongService) GetSongs(ctx context.Context, search, sortBy, order string, page, limit int) ([]dtos.SongExtendedDTO, error) {
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
		return []dtos.SongExtendedDTO{}, nil
	}

	songDTOs := make([]dtos.SongExtendedDTO, len(songs))
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

		authorsDTO := make([]dtos.AuthorDTO, 0)
		for _, songAuthor := range song.Authors {
			authorsDTO = append(authorsDTO, dtos.AuthorDTO{
				ID:   songAuthor.Author.ID,
				Name: songAuthor.Author.Name,
				Role: songAuthor.Role,
			})
		}

		tagsDTO := make([]dtos.TagDTO, 0)
		for _, songTag := range song.SongTags {
			tagsDTO = append(tagsDTO, dtos.TagDTO{
				Name: songTag.Tag.Name,
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

		songDTOs[i] = dtos.SongExtendedDTO{
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
			CreatedAt: song.CreatedAt,
			Genre:     song.Genre,
			Authors:   authorsDTO,
			SongTags:  tagsDTO,
			Comments:  []dtos.CommentDTO{},
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
	song, err := ss.GetByID(
		ctx,
		id,
		"Authors",
		"Authors.Author",
		"SongTags",
		"SongTags.Tag",
		"Comments",
		"Comments.User",
		"User",
		"Reactions",
	)
	if err != nil {
		return nil, err
	}

	songUserDTO := dtos.UserDTO{
		ID:             song.User.ID,
		Username:       song.User.Username,
		Role:           song.User.Role,
		ProfilePicture: song.User.ProfilePicture,
		ProfileInfo:    song.User.ProfileInfo,
		Followers:      int64(len(song.User.Followers)),
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
