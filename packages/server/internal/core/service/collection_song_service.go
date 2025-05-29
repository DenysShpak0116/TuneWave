package service

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type CollectionSongService struct {
	GenericService[models.CollectionSong]
}

func NewCollectionSongService(repo port.Repository[models.CollectionSong]) services.CollectionSongService {
	return &CollectionSongService{
		GenericService: GenericService[models.CollectionSong]{
			Repository: repo,
		},
	}
}
