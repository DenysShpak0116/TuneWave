package song

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SongToCollectionRequest struct {
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
// @Param        body body      SongToCollectionRequest true "Collection ID"
// @Router       /songs/{id}/add-to-collection [post]
func (sh *SongHandler) AddToCollection(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		http.Error(w, "Song ID is required", http.StatusBadRequest)
		return
	}

	songUUID, err := uuid.Parse(songID)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var request SongToCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collectionUUID, err := uuid.Parse(request.CollectionID)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	err = sh.SongService.AddToCollection(r.Context(), songUUID, collectionUUID)
	if err != nil {
		http.Error(w, "Failed to add song to collection", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Song added to collection successfully"})
}
