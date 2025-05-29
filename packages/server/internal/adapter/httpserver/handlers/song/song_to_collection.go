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

type SongCollectionRequest struct {
	CollectionID string `json:"collectionId"`
}

// AddToCollection godoc
// @Summary      Add song to collection
// @Description  Add song to collection
// @Tags         songs
// @Security    BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Song ID"
// @Param        body body      SongCollectionRequest true "Collection ID"
// @Router       /songs/{id}/add-to-collection [post]
func (sh *SongHandler) AddToCollection(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Song ID is required", nil)
		return
	}

	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
		return
	}

	var request SongCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	collectionUUID, err := uuid.Parse(request.CollectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	err = sh.SongService.AddToCollection(r.Context(), songUUID, collectionUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to add song to collection", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Song added to collection successfully"})
}

// RemoveFromCollection godoc
// @Summary      Remove song from collection
// @Description  Remove song from collection
// @Tags         songs
// @Security    BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Song ID"
// @Param        body body      SongCollectionRequest true "Collection ID"
// @Router       /songs/{id}/remove-from-collection [delete]
func (sh *SongHandler) RemoveFromCollection(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Song ID is required", nil)
		return
	}
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
		return
	}

	var request SongCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	collectionUUID, err := uuid.Parse(request.CollectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	collectionSongs, err := sh.CollectionSongService.Where(context.Background(), &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to find song in collection", err)
		return
	}

	if len(collectionSongs) == 0 {
		handlers.RespondWithError(w, r, http.StatusNotFound, "Song not found in collection", nil)
		return
	}

	collectionSong := collectionSongs[0]

	err = sh.CollectionSongService.Delete(context.Background(), collectionSong.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to remove song from collection", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Song removed from collection successfully"})
}
