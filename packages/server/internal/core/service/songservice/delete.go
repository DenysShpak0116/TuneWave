package songservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

func (ss *SongService) Delete(ctx context.Context, songID uuid.UUID) error {
	song, err := ss.getSongByID(ctx, songID)
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

	return nil
}

// TODO: replace with GetSongBYID
func (ss *SongService) getSongByID(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	songs, err := ss.Repository.NewQuery(ctx).Where("id = ?", id).Find()
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return nil, fmt.Errorf("song not found")
	}
	return &songs[0], nil
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
	songKey := extractS3Key(song.SongURL)
	if err := ss.FileStorage.Remove(ctx, songKey); err != nil {
		return fmt.Errorf("failed to delete song file from storage: %w", err)
	}

	coverKey := extractS3Key(song.CoverURL)
	if err := ss.FileStorage.Remove(ctx, coverKey); err != nil {
		return fmt.Errorf("failed to delete cover file from storage: %w", err)
	}
	return nil
}

func extractS3Key(fullURL string) string {
	const baseURL = "https://tunewavebucket.s3.eu-west-3.amazonaws.com/"
	return strings.TrimPrefix(fullURL, baseURL)
}
