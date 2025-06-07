package dto

import (
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type CollectionDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	CoverURL    string    `json:"coverUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`

	User UserDTO `json:"user"`

	CollectionSongs []SongExtendedDTO `json:"collectionSongs"`
}

func NewCollectionDTO(collection *models.Collection) *CollectionDTO {
	collectionSongs := make([]dtos.SongExtendedDTO, len(collection.CollectionSongs))
	for i, csong := range collection.CollectionSongs {
		song := csong.Song

		likes, err := ch.UserReactionService.CountWhere(ctx, &models.UserReaction{
			Type:   "like",
			SongID: song.ID,
		})
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting likes", err)
			return
		}
		dislikes, err := ch.UserReactionService.CountWhere(ctx, &models.UserReaction{
			Type:   "dislike",
			SongID: song.ID,
		})
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting dislikes", err)
			return
		}

		collectionSongs[i] = dtos.SongExtendedDTO{
			ID:         song.ID,
			Title:      song.Title,
			Duration:   formatDuration(time.Duration(song.Duration)),
			CoverURL:   song.CoverURL,
			Listenings: song.Listenings,
			User: dtos.UserDTO{
				ID:             song.User.ID,
				Username:       song.User.Username,
				ProfilePicture: song.User.ProfilePicture,
				ProfileInfo:    song.User.ProfileInfo,
				Followers:      int64(len(song.User.Followers)),
			},
			CreatedAt: song.CreatedAt,
			Genre:     song.Genre,
			SongURL:   song.SongURL,
			Likes:     likes,
			Dislikes:  dislikes,
			// Authors, Tags, Comments можно подтянуть отдельно, если нужно
		}
	}

	dto := &dtos.CollectionExtendedDTO{
		ID:          collection.ID,
		Title:       collection.Title,
		Description: collection.Description,
		CoverURL:    collection.CoverURL,
		CreatedAt:   collection.CreatedAt,
		User: dtos.UserDTO{
			ID:             collection.User.ID,
			Username:       collection.User.Username,
			ProfilePicture: collection.User.ProfilePicture,
			ProfileInfo:    collection.User.ProfileInfo,
			Followers:      int64(len(collection.User.Followers)), // Аналогично
		},
		CollectionSongs: collectionSongs,
	}
}
