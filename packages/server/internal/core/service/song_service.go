package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type SongService struct {
	*GenericService[models.Song]
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
		GenericService:           NewGenericService(songRepo),
		FileStorage:              fileStorage,
		AuthorRepository:         authorRepo,
		SongAuthorRepository:     songAuthorRepo,
		TagsRepository:           tagRepository,
		SongTagsRepository:       songTagRepository,
		ReactionsRepository:      reactionsRepository,
		CollectionSongRepository: collectionSongReporitory,
	}
}

func (ss *SongService) GetSongs(ctx context.Context, params services.SearchSongsParams, preloads ...string) ([]models.Song, error) {
	query := ss.Repository.NewQuery(ctx).
		Join("LEFT JOIN song_authors ON song_authors.song_id = songs.id").
		Join("LEFT JOIN authors ON authors.id = song_authors.author_id").
		Where("songs.title ILIKE ? OR authors.name ILIKE ? OR LOWER(songs.genre) = LOWER(?)",
			"%"+params.Search+"%", "%"+params.Search+"%", params.Search).
		Order(fmt.Sprintf("songs.%s %s", params.SortBy, params.Order)).
		Group("songs.id").
		Skip((params.Page - 1) * params.Limit).
		Take(params.Limit).
		Preload(preloads...)

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

func (ss *SongService) UpdateSong(ctx context.Context, songParams services.UpdateSongParams) error {
	songs, err := ss.Repository.NewQuery(ctx).
		Where("id = ?", songParams.SongID).
		Find()
	if err != nil {
		return fmt.Errorf("failed to find song: %w", err)
	}
	if len(songs) == 0 {
		return fmt.Errorf("song not found")
	}
	song := &songs[0]

	if songParams.Title != "" {
		song.Title = songParams.Title
	}
	if songParams.Genre != "" {
		song.Genre = songParams.Genre
	}

	if songParams.Song != nil && songParams.SongHeader != nil {
		oldSongKey := helpers.ExtractS3Key(song.SongURL)
		if err := ss.FileStorage.Remove(ctx, oldSongKey); err != nil {
			return fmt.Errorf("failed to remove old song file: %w", err)
		}

		key := fmt.Sprintf("music/%s/%d-%s", song.UserID, time.Now().Unix(), songParams.SongHeader.Filename)
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, songParams.Song); err != nil {
			return err
		}

		url, err := ss.FileStorage.Save(ctx, key, buf)
		if err != nil {
			return err
		}

		duration, err := helpers.GetAudioDuration(&readSeekCloser{bytes.NewReader(buf.Bytes())})
		if err != nil {
			return fmt.Errorf("failed to get audio duration: %w", err)
		}

		song.SongURL = url
		song.Duration = duration
	}

	if songParams.Cover != nil && songParams.CoverHeader != nil {
		oldCoverKey := helpers.ExtractS3Key(song.SongURL)
		if err := ss.FileStorage.Remove(ctx, oldCoverKey); err != nil {
			return fmt.Errorf("failed to remove old cover file: %w", err)
		}

		key := fmt.Sprintf("covers/%s/%d-%s", song.UserID, time.Now().Unix(), songParams.CoverHeader.Filename)
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, songParams.Cover); err != nil {
			return err
		}

		url, err := ss.FileStorage.Save(ctx, key, buf)
		if err != nil {
			return err
		}

		song.CoverURL = url
	}

	if err := ss.Repository.Update(ctx, song); err != nil {
		return err
	}

	if songParams.Artists != nil {
		err := ss.SongAuthorRepository.NewQuery(ctx).
			Where("song_id = ?", song.ID).
			Delete()
		if err != nil {
			return fmt.Errorf("failed to delete old song authors: %w", err)
		}

		if err := ss.associateAuthors(ctx, song, songParams.Artists); err != nil {
			return err
		}
	}

	if songParams.Tags != nil {
		err := ss.SongTagsRepository.NewQuery(ctx).
			Where("song_id = ?", song.ID).
			Delete()
		if err != nil {
			return fmt.Errorf("failed to delete old song tags: %w", err)
		}

		if err := ss.associateTags(ctx, song, songParams.Tags); err != nil {
			return err
		}
	}

	return nil
}

func (ss *SongService) SaveSong(ctx context.Context, songParams services.SaveSongParams) (*models.Song, error) {
	songURL, duration, err := ss.saveSongFile(ctx, songParams)
	if err != nil {
		return nil, err
	}

	coverURL, err := ss.saveCoverFile(ctx, songParams)
	if err != nil {
		return nil, err
	}

	song := &models.Song{
		Title:      songParams.Title,
		Genre:      songParams.Genre,
		SongURL:    songURL,
		CoverURL:   coverURL,
		Duration:   duration,
		Listenings: 0,
		UserID:     songParams.UserID,
	}

	if err := ss.Repository.Add(ctx, song); err != nil {
		return nil, err
	}

	if err := ss.associateAuthors(ctx, song, songParams.Artists); err != nil {
		return nil, err
	}

	if err := ss.associateTags(ctx, song, songParams.Tags); err != nil {
		return nil, err
	}

	return song, nil
}

func (ss *SongService) saveSongFile(ctx context.Context, songParams services.SaveSongParams) (string, time.Duration, error) {
	key := fmt.Sprintf("music/%s/%d-%s", songParams.UserID, time.Now().Unix(), songParams.SongHeader.Filename)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, songParams.Song); err != nil {
		return "", 0, err
	}

	url, err := ss.FileStorage.Save(ctx, key, buf)
	if err != nil {
		return "", 0, err
	}

	duration, err := helpers.GetAudioDuration(&readSeekCloser{bytes.NewReader(buf.Bytes())})
	fmt.Printf("Duration: %v\n", duration)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get audio duration: %w", err)
	}

	return url, duration, nil
}

func (ss *SongService) saveCoverFile(ctx context.Context, songParams services.SaveSongParams) (string, error) {
	key := fmt.Sprintf("covers/%s/%d-%s", songParams.UserID, time.Now().Unix(), songParams.CoverHeader.Filename)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, songParams.Cover); err != nil {
		return "", err
	}

	url, err := ss.FileStorage.Save(ctx, key, buf)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (ss *SongService) associateAuthors(ctx context.Context, song *models.Song, artistNames []string) error {
	for _, name := range artistNames {
		var author *models.Author

		found, err := ss.AuthorRepository.NewQuery(ctx).Where("name = ?", name).Take(1).Find()
		if err != nil {
			return fmt.Errorf("failed to check existing author: %w", err)
		}

		if len(found) > 0 {
			author = &found[0]
		} else {
			author = &models.Author{Name: name}
			if err := ss.AuthorRepository.Add(ctx, author); err != nil {
				return fmt.Errorf("failed to create new author: %w", err)
			}
		}

		songAuthor := &models.SongAuthor{
			SongID:   song.ID,
			AuthorID: author.ID,
			Role:     "primary",
		}
		if err := ss.SongAuthorRepository.Add(ctx, songAuthor); err != nil {
			return fmt.Errorf("failed to associate author with song: %w", err)
		}
	}
	return nil
}

func (ss *SongService) associateTags(ctx context.Context, song *models.Song, tagNames []string) error {
	for _, name := range tagNames {
		var tag *models.Tag

		found, err := ss.TagsRepository.NewQuery(ctx).Where("name = ?", name).Take(1).Find()
		if err != nil {
			return fmt.Errorf("failed to check existing tag: %w", err)
		}

		if len(found) > 0 {
			tag = &found[0]
		} else {
			tag = &models.Tag{Name: name}
			if err := ss.TagsRepository.Add(ctx, tag); err != nil {
				return fmt.Errorf("failed to create new tag: %w", err)
			}
		}

		songTag := &models.SongTag{
			SongID: song.ID,
			TagID:  tag.ID,
		}
		if err := ss.SongTagsRepository.Add(ctx, songTag); err != nil {
			return fmt.Errorf("failed to associate tag with song: %w", err)
		}
	}
	return nil
}

func (ss *SongService) SetReaction(ctx context.Context, songID uuid.UUID, userID uuid.UUID, reactionType string) (int64, int64, error) {
	reactions, err := ss.ReactionsRepository.NewQuery(ctx).
		Where("song_id = ? AND user_id = ?", songID, userID).
		Find()
	if err != nil {
		return 0, 0, err
	}

	if len(reactions) == 0 {
		reaction := models.UserReaction{
			SongID: songID,
			UserID: userID,
			Type:   reactionType,
		}
		if err := ss.ReactionsRepository.Add(ctx, &reaction); err != nil {
			return 0, 0, err
		}
	} else {
		existingReaction := reactions[0]
		if existingReaction.Type == reactionType {
			if err := ss.ReactionsRepository.Delete(ctx, existingReaction.ID); err != nil {
				return 0, 0, err
			}
		} else {
			existingReaction.Type = reactionType
			if err := ss.ReactionsRepository.Update(ctx, &existingReaction); err != nil {
				return 0, 0, err
			}
		}
	}

	dislikes, err := ss.ReactionsCount(ctx, songID, "dislike")
	if err != nil {
		return 0, 0, err
	}
	likes, err := ss.ReactionsCount(ctx, songID, "like")
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

func (ss *SongService) IsReactedByUser(ctx context.Context, songID uuid.UUID, userID uuid.UUID) (string, error) {
	reactions, err := ss.ReactionsRepository.NewQuery(ctx).
		Where("song_id = ? AND user_id = ?", songID, userID).
		Find()
	if err != nil {
		return "none", err
	}
	if len(reactions) == 0 {
		return "none", nil
	}

	reactionType := reactions[0].Type
	return reactionType, nil
}

func (ss *SongService) GetGenres(ctx context.Context) []string {
	return ss.Repository.Distinct(ctx, "genre")
}

func (ss *SongService) GetGenresMostPopularSong(ctx context.Context, genre string) (*models.Song, error) {
	songs, err := ss.Repository.NewQuery(ctx).
		Join("JOIN user_reactions r ON r.song_id = songs.id").
		Where("genre = ? AND r.type = ?", genre, "like").
		Group("songs.id").
		Order("COUNT(r.id) DESC").
		Take(1).
		Find()
	if err != nil {
		return nil, err
	}

	if len(songs) == 0 {
		return nil, nil
	}

	return &songs[0], nil
}

func (ss *SongService) Delete(ctx context.Context, songIDs ...uuid.UUID) error {
	for _, songID := range songIDs {
		song, err := ss.GetByID(ctx, songID)
		if err != nil {
			return err
		}

		if err := ss.cleanupSongAuthors(ctx, songID); err != nil {
			return err
		}

		if err := ss.cleanupSongTags(ctx, songID); err != nil {
			return err
		}

		if err := ss.Repository.Delete(ctx, song.ID); err != nil {
			return fmt.Errorf("failed to delete song: %w", err)
		}

		if err := ss.removeSongFiles(ctx, song); err != nil {
			return err
		}
	}

	return nil
}

func (ss *SongService) cleanupSongAuthors(ctx context.Context, songID uuid.UUID) error {
	songAuthors, err := ss.SongAuthorRepository.NewQuery(ctx).Where("song_id = ?", songID).Find()
	if err != nil {
		return fmt.Errorf("failed to fetch song authors: %w", err)
	}

	if err := ss.SongAuthorRepository.NewQuery(ctx).Where("song_id = ?", songID).Delete(); err != nil {
		return fmt.Errorf("failed to delete song_author entries: %w", err)
	}

	for _, sa := range songAuthors {
		count, err := ss.SongAuthorRepository.NewQuery(ctx).Where("author_id = ?", sa.AuthorID).Count()
		if err != nil {
			return fmt.Errorf("failed to count songs for author: %w", err)
		}
		if count == 0 {
			if err := ss.AuthorRepository.Delete(ctx, sa.AuthorID); err != nil {
				return fmt.Errorf("failed to delete author: %w", err)
			}
		}
	}

	return nil
}

func (ss *SongService) cleanupSongTags(ctx context.Context, songID uuid.UUID) error {
	songTags, err := ss.SongTagsRepository.NewQuery(ctx).Where("song_id = ?", songID).Find()
	if err != nil {
		return fmt.Errorf("failed to fetch song tags: %w", err)
	}

	if err := ss.SongTagsRepository.NewQuery(ctx).Where("song_id = ?", songID).Delete(); err != nil {
		return fmt.Errorf("failed to delete song_tags entries: %w", err)
	}

	for _, st := range songTags {
		count, err := ss.SongTagsRepository.NewQuery(ctx).Where("tag_id = ?", st.TagID).Count()
		if err != nil {
			return fmt.Errorf("failed to count songs for tag: %w", err)
		}
		if count == 0 {
			if err := ss.TagsRepository.Delete(ctx, st.TagID); err != nil {
				return fmt.Errorf("failed to delete tag: %w", err)
			}
		}
	}

	return nil
}

func (ss *SongService) removeSongFiles(ctx context.Context, song *models.Song) error {
	songKey := helpers.ExtractS3Key(song.SongURL)
	if err := ss.FileStorage.Remove(ctx, songKey); err != nil {
		return fmt.Errorf("failed to delete song file from storage: %w", err)
	}

	coverKey := helpers.ExtractS3Key(song.CoverURL)
	if err := ss.FileStorage.Remove(ctx, coverKey); err != nil {
		return fmt.Errorf("failed to delete cover file from storage: %w", err)
	}
	return nil
}

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
