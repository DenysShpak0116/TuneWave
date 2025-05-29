package songservice

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

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
