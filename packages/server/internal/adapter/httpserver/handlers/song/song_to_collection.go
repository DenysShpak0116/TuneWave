package song

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SongCollectionRequest struct {
	CollectionID string `json:"collectionId"`
}

// TODO: Handle only collection creator can add songs to it
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
	const op = "adapter.httpserver.handlers.SongHandler.AddToCollection"
	logger := sh.logger.With(
		slog.String("op", op),
	)

	ctx := r.Context()
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.BadRequest("invalid user ID")
	}

	var request SongCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.BadRequest("invalid request body")
	}
	collectionUUID, err := uuid.Parse(request.CollectionID)
	if err != nil {
		return helpers.BadRequest("invalid collection ID")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid song ID")
	}

	collectionSong := &models.CollectionSong{
		SongID:       songUUID,
		CollectionID: collectionUUID,
	}
	if _, err := sh.collectionSongService.First(ctx, collectionSong); !errors.Is(err, service.ErrNotFound) {
		if err != nil {
			return helpers.InternalServerError("could not add song to collection")
		}

		return helpers.NewAPIError(http.StatusConflict, "this song is already in collection")
	}

	if err := sh.songService.AddToCollection(ctx, songUUID, collectionUUID); err != nil {
		return helpers.InternalServerError("failed to add song to collection")
	}

	event := &models.Event{
		BaseModel:  models.BaseModel{},
		UserID:     userUUID,
		EventType:  models.AddToPlaylistEvent,
		TrackID:    &songUUID,
		PlaylistID: &collectionUUID,
		Metadata: datatypes.JSONMap{
			"source": "add_to_collection_endpoint",
		},
	}
	if err := sh.eventService.Create(ctx, event); err != nil {
		logger.Error("failed to create add_to_playlist event", "err", err.Error())
	}

	render.JSON(w, r, map[string]string{"message": "Song added to collection"})
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
	ctx := r.Context()
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.BadRequest("invalid user ID")
	}

	var request SongCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.BadRequest("invalid request body")
	}
	collectionUUID, err := uuid.Parse(request.CollectionID)
	if err != nil {
		return helpers.BadRequest("invalid collection ID")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid song ID")
	}

	collectionSong, err := sh.collectionSongService.First(ctx, &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	})
	if err != nil {
		return helpers.InternalServerError("failed to find song in collection")
	}

	if err = sh.collectionSongService.Delete(ctx, collectionSong.ID); err != nil {
		return helpers.InternalServerError("failed to remove song from collection")
	}

	event := &models.Event{
		BaseModel:  models.BaseModel{},
		UserID:     userUUID,
		EventType:  models.RemoveFromPlaylistEvent,
		TrackID:    &songUUID,
		PlaylistID: &collectionUUID,
		Metadata: datatypes.JSONMap{
			"source": "remove_from_collection_endpoint",
		},
	}
	if err := sh.eventService.Create(ctx, event); err != nil {
		fmt.Println("failed to create remove_from_playlist event", "err", err.Error())
	}
	render.JSON(w, r, map[string]string{"message": "Song removed from collection"})
	return nil
}
