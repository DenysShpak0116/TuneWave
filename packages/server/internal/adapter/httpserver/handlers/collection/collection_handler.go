package collection

import (
	"context"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CollectionHandler struct {
	CollectionService     services.CollectionService
	UserCollectionService services.UserCollectionService
}

func NewCollectionHandler(
	collectionService services.CollectionService,
	userCollectionService services.UserCollectionService,
) *CollectionHandler {
	return &CollectionHandler{
		CollectionService:     collectionService,
		UserCollectionService: userCollectionService,
	}
}

// CreateCollection godoc
// @Summary Create a new collection
// @Description Creates a new collection. Returns the created collection object.
// @Tags collections
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
	userUUID, err := uuid.Parse(userID)
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
			UserID:      userUUID,
		},
	)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error saving collection", err)
		return
	}

	userCollection := &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collection.ID,
	}

	if err := ch.UserCollectionService.Create(context.Background(), userCollection); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error creating user collection", err)
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
// @Tags collections
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
// @Tags collections
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
// @Tags collections
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

// GetUsersCollections godoc
// @Summary      Get user's collections
// @Description  Retrieves all collections belonging to the authenticated user
// @Tags         collections
// @Security     BearerAuth
// @Produce      json
// @Router       /collections/users-collections [get]
func (ch *CollectionHandler) GetUsersCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r.Context())
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}

	userCollections, err := ch.UserCollectionService.Where(context.Background(), &models.UserCollection{
		UserID: userUUID,
	}, "Collection")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting user collections", err)
		return
	}

	if len(userCollections) == 0 {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, []dtos.UsersCollectionDTO{})
		return
	}

	collections := make([]models.Collection, len(userCollections))
	for i, userCollection := range userCollections {
		collections[i] = userCollection.Collection
	}

	var usersCollectionsDTOs []dtos.UsersCollectionDTO
	if err := dtoMapper.Map(&usersCollectionsDTOs, collections); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Cannot get collections", err)
		return
	}

	render.JSON(w, r, usersCollectionsDTOs)
}

// GetCollections godoc
// @Summary Get all collections
// @Description Get all collections. Returns a list of collections.
// @Tags collections
// @Security     BearerAuth
// @Produce  json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort query string false "Sort by (title, created_at)"
// @Param order query string false "Order (asc, desc)"
// @Router /collections [get]
func (ch *CollectionHandler) GetCollections(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	pageParam := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	if sort != "title" && sort != "created_at" {
		sort = "created_at"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	var limit int
	if limitParam == "" {
		limit = 10
	} else {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid limit", err)
			return
		}
	}
	if limit <= 0 {
		limit = 10
	}

	var page int
	if pageParam == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid page", err)
			return
		}
	}
	if page <= 0 {
		page = 1
	}

	collections, err := ch.CollectionService.GetMany(context.Background(), limit, page, sort, order)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting collections", err)
		return
	}

	var collectionsDTOs []dtos.CollectionDTO
	if err := dtoMapper.Map(&collectionsDTOs, collections); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Cannot get collections", err)
		return
	}

	render.JSON(w, r, collectionsDTOs)
}

// AddCollectionToUser godoc
// @Summary Add a collection to a user
// @Description Adds a collection to a user. Returns the updated list of collections.
// @Tags collections
// @Security     BearerAuth
// @Produce  json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/add-to-user [post]
func (ch *CollectionHandler) AddCollectionToUser(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")

	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	userID, err := helpers.GetUserID(r.Context())
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}

	if userCollections, err := ch.UserCollectionService.Where(context.Background(), &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collectionUUID,
	}); err != nil || len(userCollections) > 0 {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "could not add collection", err)
		return
	}

	err = ch.UserCollectionService.Create(context.Background(), &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collectionUUID,
	})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error adding song to collection", err)
		return
	}

	userCollections, err := ch.UserCollectionService.Where(context.Background(), &models.UserCollection{
		UserID: userUUID,
	}, "Collection")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting user collections", err)
		return
	}

	if len(userCollections) == 0 {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, []dtos.UsersCollectionDTO{})
		return
	}

	collections := make([]models.Collection, len(userCollections))
	for i, userCollection := range userCollections {
		collections[i] = userCollection.Collection
	}

	var usersCollectionsDTOs []dtos.UsersCollectionDTO
	if err := dtoMapper.Map(&usersCollectionsDTOs, collections); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Cannot get collections", err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, usersCollectionsDTOs)
}

// RemoveCollectionFromUser godoc
// @Summary Remove a collection from a user
// @Description Removes a collection from a user. Returns no content on success.
// @Tags collections
// @Security     BearerAuth
// @Produce  json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/remove-from-user [delete]
func (ch *CollectionHandler) RemoveCollectionFromUser(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")

	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	userID, err := helpers.GetUserID(r.Context())
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}

	userCollection, err := ch.UserCollectionService.Where(context.Background(), &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collectionUUID,
	})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error getting user collection", err)
		return
	}
	if len(userCollection) == 0 {
		handlers.RespondWithError(w, r, http.StatusNotFound, "User collection not found", nil)
		return
	}

	userCollectionID := userCollection[0].ID

	err = ch.UserCollectionService.Delete(context.Background(), userCollectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error removing song from collection", err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// GetCollectionSongs godoc
// @Tags collections
// @Produce json
// @Param id path string true "Collection id"
// @Param search query string false "Search by title"
// @Param sortBy query string false "Sort by field (added_at, title)" default(added_at)
// @Param order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Router /collections/{id}/songs [get]
func (ch *CollectionHandler) GetCollectionSongs(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Wrong collectionID", err)
		return
	}

	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if sortBy == "" {
		sortBy = "created_at"
	}

	if order == "" {
		order = "desc"
	}

	collectionSongs, err := ch.CollectionService.GetCollectionSongs(
		context.Background(),
		collectionUUID,
		search,
		sortBy,
		order,
		page,
		limit,
	)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Coult not retrieve collection songs", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, collectionSongs)
}
