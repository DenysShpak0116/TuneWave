package collection

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CollectionHandler struct {
	CollectionService services.CollectionService
}

func NewCollectionHandler(collectionService services.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		CollectionService: collectionService,
	}
}

// CreateCollection godoc
// @Summary Create a new collection
// @Description Creates a new collection. Returns the created collection object.
// @Tags Collections
// @Security     BearerAuth
// @Accept  multipart/form-data
// @Produce  json
// @Param userId formData string true "User ID"
// @Param title formData string true "Collection title"
// @Param description formData string true "Collection description"
// @Param cover formData file true "Collection cover image"
// @Router /collections [post]
func (ch *CollectionHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Error parsing form", err)
		return
	}

	userID := r.FormValue("userId")
	userIDuuid, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	coverFile, coverHeader, err := r.FormFile("cover")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Error getting cover file", err)
		return
	}

	collection, err := ch.CollectionService.SaveCollection(
		context.Background(), services.SaveCollectionParams{
			Title:       title,
			Description: description,
			CoverHeader: coverHeader,
			Cover:       coverFile,
			UserID:      userIDuuid,
		},
	)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error saving collection", err)
		return
	}

	collectionDTO, err := ch.CollectionService.GetFullDTOByID(context.Background(), collection.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting collection", err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, collectionDTO)
}

// GetCollectionByID godoc
// @Summary Get collection by ID
// @Description Get collection by ID. Returns the collection object.
// @Tags Collections
// @Security     BearerAuth
// @Produce  json
// @Param id path string true "Collection ID"
// @Router /collections/{id} [get]
func (ch *CollectionHandler) GetCollectionByID(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	collectionDTO, err := ch.CollectionService.GetFullDTOByID(r.Context(), collectionUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting collection", err)
		return
	}

	render.JSON(w, r, collectionDTO)
}

// DeleteCollection godoc
// @Summary Delete a collection
// @Description Deletes a collection by its ID. Returns no content on success.
// @Tags Collections
// @Security     BearerAuth
// @Param id path string true "Collection ID"
// @Router /collections/{id} [delete]
func (ch *CollectionHandler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	err = ch.CollectionService.Delete(r.Context(), collectionUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error deleting collection", err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// UpdateCollection godoc
// @Summary Update a collection
// @Description Updates a collection by its ID. Returns the updated collection object.
// @Tags Collections
// @Security     BearerAuth
// @Accept  multipart/form-data
// @Produce  json
// @Param id path string true "Collection ID"
// @Param title formData string true "Collection title"
// @Param description formData string true "Collection description"
// @Param cover formData file true "Collection cover image"
// @Router /collections/{id} [put]
func (ch *CollectionHandler) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Error parsing form", err)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	var coverFile multipart.File
	var coverHeader *multipart.FileHeader

	if r.MultipartForm != nil && len(r.MultipartForm.File["cover"]) > 0 {
		coverFile, coverHeader, err = r.FormFile("cover")
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusBadRequest, "Error getting cover file", err)
			return
		}
	}
	prevCollection, err := ch.CollectionService.GetFullDTOByID(r.Context(), collectionUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting collection", err)
		return
	}

	_, err = ch.CollectionService.UpdateCollection(
		r.Context(), collectionUUID, services.UpdateCollectionParams{
			UserID:      prevCollection.User.ID,
			Title:       title,
			Description: description,
			CoverHeader: coverHeader,
			Cover:       coverFile,
		},
	)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error updating collection", err)
		return
	}

	newCollection, err := ch.CollectionService.GetFullDTOByID(r.Context(), prevCollection.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting collection", err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, newCollection)
}
