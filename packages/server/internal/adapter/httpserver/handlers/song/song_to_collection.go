package song

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SongCollectionRequest struct {
	CollectionID string `json:"collectionId"`
}

// AddToCollection godoc
// @Summary        Add song to collection
// @Description    Add song to collection
// @Tags           songs
// @Security       BearerAuth
// @Accept         json
// @Produce        json
// @Param          id   path      string  true  "Song ID"
// @Param          body body      SongCollectionRequest true "Collection ID"
// @Router         /songs/{id}/add-to-collection [post]
func (sh *SongHandler) AddToCollection(w http.ResponseWriter, r *http.Request) error {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		return helpers.NewAPIError(http.StatusBadRequest, "song ID is requiredr")
	}

	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	var request SongCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request body")
	}

	collectionUUID, err := uuid.Parse(request.CollectionID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	err = sh.SongService.AddToCollection(r.Context(), songUUID, collectionUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to add song to collection")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Song added to collection successfully"})
	return nil
}

// RemoveFromCollection godoc
// @Summary      Remove song from collection
// @Description  Remove song from collection
// @Tags         songs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string  true  "Song ID"
// @Param        body body SongCollectionRequest true "Collection ID"
// @Router       /songs/{id}/remove-from-collection [delete]
func (sh *SongHandler) RemoveFromCollection(w http.ResponseWriter, r *http.Request) error {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		return helpers.NewAPIError(http.StatusBadRequest, "song ID is required")
	}
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	var request SongCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request body")
	}
	collectionUUID, err := uuid.Parse(request.CollectionID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	collectionSongs, err := sh.CollectionSongService.Where(context.Background(), &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to find song in collection")
	}

	if len(collectionSongs) == 0 {
		return helpers.NewAPIError(http.StatusNotFound, "song not found in collection")
	}

	collectionSong := collectionSongs[0]

	err = sh.CollectionSongService.Delete(context.Background(), collectionSong.ID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to remove song from collection")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Song removed from collection successfully"})
	return nil
}
