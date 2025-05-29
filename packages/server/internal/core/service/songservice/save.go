package songservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

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
