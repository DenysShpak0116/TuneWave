//go:generate mockgen -source=collection_song_service.go -destination=../../service/mocks/collection_song_service_mock.go -package=mocks -typed

package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type CollectionSongService interface {
	Service[models.CollectionSong]
}
