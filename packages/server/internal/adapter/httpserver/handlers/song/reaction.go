package song

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type songReactionRequest struct {
	ReactionType string `json:"reactionType"`
	UserID       string `json:"userId"`
}

// SetReaction godoc
// @Summary      Set reaction to song
// @Description  Set reaction to song
// @Security     BearerAuth
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param id path string true "Song ID"
// @Param body body songReactionRequest true "Reaction request body"
// @Router /songs/{id}/reaction [post]
func (sh *SongHandler) SetReaction(w http.ResponseWriter, r *http.Request) error {
	songUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	var request songReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request body")
	}

	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "something wrong with userID")
	}

	likes, dislikes, err := sh.songService.SetReaction(r.Context(), songUUID, userUUID, request.ReactionType)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to set reaction")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"likes":    likes,
		"dislikes": dislikes,
	})

	return nil
}

// CheckReaction godoc
// @Summary      Check reaction to song
// @Description  Check reaction to song
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param id path string true "Song ID"
// @Param userId path string true "User ID"
// @Router /songs/{id}/is-reacted/{userId} [get]
func (sh *SongHandler) CheckReaction(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userId")
	if userID == "undefined" {
		render.JSON(w, r, map[string]string{
			"type": "none",
		})
		return nil
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user id")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song id")
	}

	reactionType, err := sh.songService.IsReactedByUser(r.Context(), songUUID, userUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to check reaction")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"type": reactionType,
	})
	return nil
}

// ListenSong godoc
// @Summary      Add listening to song
// @Description  Add listening to song
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param id path string true "Song ID"
// @Param userId path string true "User ID"
// @Router /songs/{id}/listen/{userId} [post]
func (sh *SongHandler) ListenSong(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := chi.URLParam(r, "userId")
	if userID == "undefined" {
		render.NoContent(w, r)
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "song id is wrong")
	}

	song, err := sh.songService.GetByID(ctx, songUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "song does not exist")
	}
	if err := sh.songService.Update(ctx, &models.Song{
		BaseModel: models.BaseModel{
			ID: song.ID,
		},
		Listenings: song.Listenings + 1,
	}); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "v")
	}

	render.NoContent(w, r)
	return nil
}
