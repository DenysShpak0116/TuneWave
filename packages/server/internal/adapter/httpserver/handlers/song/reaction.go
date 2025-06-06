package song

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
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
func (sh *SongHandler) SetReaction(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid song id", err)
		return
	}

	var request songReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid request body", err)
		return
	}

	userID := request.UserID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "something wrong with userID", err)
		return
	}
	likes, dislikes, err := sh.SongService.SetReaction(r.Context(), songUUID, userUUID, request.ReactionType)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to set reaction", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"likes":    likes,
		"dislikes": dislikes,
	})
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
func (sh *SongHandler) CheckReaction(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if userID == "undefined" {
		render.JSON(w, r, map[string]string{
			"type": "none",
		})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid user id", err)
		return
	}

	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid song id", err)
		return
	}

	reactionType, err := sh.SongService.IsReactedByUser(r.Context(), songUUID, userUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to check reaction", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"type": reactionType,
	})
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
func (sh *SongHandler) ListenSong(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if userID == "undefined" {
		render.Status(r, http.StatusNoContent)
	}

	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "song id is wrong", err)
		return
	}

	song, err := sh.SongService.GetByID(
		context.Background(),
		songUUID,
		"Authors",
		"Authors.Author",
		"SongTags",
		"SongTags.Tag",
		"Comments",
		"Comments.User",
		"User",
		"Reactions",
	)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "song does not exist", err)
		return

	}
	if _, err := sh.SongService.Update(context.Background(), &models.Song{
		BaseModel: models.BaseModel{
			ID: song.ID,
		},
		Listenings: song.Listenings + 1,
	}); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "can not add listening to song", err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
