package services

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type SongService interface {
	Service[models.Song]
	GetSongs(
		ctx context.Context,
		search string,
		sortBy string,
		order string,
		page int,
		limit int,
	) ([]models.Song, error)
}
