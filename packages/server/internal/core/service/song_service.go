package service

import (
	"context"
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
)

type SongService struct {
	*GenericService[models.Song]
}

func NewSongService(songRepo port.Repository[models.Song]) *SongService {
	return &SongService{
		GenericService: NewGenericService(songRepo),
	}
}

func (ss *SongService) GetSongs(
	ctx context.Context,
	search, sortBy, order string,
	page, limit int,
) ([]models.Song, error) {
	songs, err := ss.Repository.NewQuery(ctx).
		Where("title LIKE ?", "%"+search+"%").
		Order(fmt.Sprintf("%s %s", sortBy, order)).
		Skip((page - 1) * limit).
		Take(limit).
		Find()
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return nil, fmt.Errorf("no songs found for search term: %s", search)
	}
	return songs, nil
}
