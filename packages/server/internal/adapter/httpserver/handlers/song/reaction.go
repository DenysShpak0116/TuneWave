package song

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"gorm.io/datatypes"
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
		return helpers.BadRequest("invalid song ID")
	}

	var request songReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.BadRequest("invalid request body")
	}

	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		return helpers.BadRequest("something wrong with userID")
	}

	likes, dislikes, reactionAction, err := sh.songService.SetReaction(r.Context(), songUUID, userUUID, request.ReactionType)
	if err != nil {
		return helpers.InternalServerError("failed to set reaction")
	}

	event := &models.Event{
		BaseModel: models.BaseModel{},
		UserID:    userUUID,
		EventType: models.LikeEvent,
		TrackID:   &songUUID,
		Metadata: datatypes.JSONMap{
			"action": reactionAction,
		},
	}
	if err := sh.eventService.Create(r.Context(), event); err != nil {
		sh.logger.Error("failed to create reaction event", "error", err)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]any{
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
		return helpers.BadRequest("invalid user id")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid song id")
	}

	reactionType, err := sh.songService.IsReactedByUser(r.Context(), songUUID, userUUID)
	if err != nil {
		return helpers.InternalServerError("failed to check reaction")
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
	const op = "adapter.httpserver.handlers.SongHandler.ListenSong"
	logger := sh.logger.With(
		slog.String("op", op),
	)

	ctx := r.Context()
	userID := chi.URLParam(r, "userId")
	if userID == "undefined" {
		render.NoContent(w, r)
		return nil
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.BadRequest("user id is wrong")
	}

	songUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("song id is wrong")
	}

	song, err := sh.songService.GetByID(ctx, songUUID)
	if err != nil {
		return helpers.BadRequest("song does not exist")
	}

	updateSongParams := &models.Song{
		BaseModel:  models.BaseModel{ID: song.ID},
		Listenings: song.Listenings + 1,
	}
	if err := sh.songService.Update(ctx, updateSongParams); err != nil {
		return helpers.InternalServerError("v")
	}

	event := &models.Event{
		UserID:    userUUID,
		EventType: models.PlayEvent,
		TrackID:   &song.ID,
		Metadata: datatypes.JSONMap{
			"source": "direct_play",
		},
	}

	if err := sh.eventService.Create(ctx, event); err != nil {
		logger.Error("failed to create play event", "err", err.Error())
	}

	render.NoContent(w, r)
	return nil
}
