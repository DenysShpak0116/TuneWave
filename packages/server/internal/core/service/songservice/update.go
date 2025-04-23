package songservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

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
		oldSongKey := extractS3Key(song.SongURL)
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
		oldCoverKey := extractS3Key(song.SongURL)
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

	if _, err := ss.Repository.Update(ctx, song); err != nil {
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
