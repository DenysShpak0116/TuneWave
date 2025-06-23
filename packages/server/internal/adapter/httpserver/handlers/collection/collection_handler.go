package collection

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/query"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CollectionHandler struct {
	collectionService     services.CollectionService
	userCollectionService services.UserCollectionService
	userReactionService   services.UserReactionService
	userService           services.UserService
	dtoBuilder            *dto.DTOBuilder
}

func NewCollectionHandler(
	collectionService services.CollectionService,
	userCollectionService services.UserCollectionService,
	userReactionService services.UserReactionService,
	userService services.UserService,
	dtoBuilder *dto.DTOBuilder,
) *CollectionHandler {
	return &CollectionHandler{
		collectionService:     collectionService,
		userCollectionService: userCollectionService,
		userReactionService:   userReactionService,
		userService:           userService,
		dtoBuilder:            dtoBuilder,
	}
}

// CreateCollection godoc
// @Summary Create a new collection
// @Description Creates a new collection. Returns the created collection object.
// @Tags collections
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Collection title"
// @Param description formData string true "Collection description"
// @Param cover formData file true "Collection cover image"
// @Router /collections [post]
func (ch *CollectionHandler) CreateCollection(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "error parsing form")
	}

	userID, _ := helpers.GetUserID(ctx)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	coverFile, coverHeader, err := r.FormFile("cover")
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "error getting cover file")
	}

	collection, err := ch.collectionService.SaveCollection(
		ctx, services.SaveCollectionParams{
			Title:       title,
			Description: description,
			CoverHeader: coverHeader,
			Cover:       coverFile,
			UserID:      userUUID,
		},
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error saving collection")
	}

	userCollection := &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collection.ID,
	}
	if err := ch.userCollectionService.Create(ctx, userCollection); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error creating user collection")
	}

	preloads := []string{"User"}
	newCollection, err := ch.collectionService.GetByID(ctx, collection.ID, preloads...)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting collection")
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ch.dtoBuilder.BuildCollectionDTO(newCollection))
	return nil
}

// GetCollectionByID godoc
// @Summary Get collection by ID
// @Description Get collection by ID. Returns the collection object.
// @Tags collections
// @Security     BearerAuth
// @Produce  json
// @Param id path string true "Collection ID"
// @Router /collections/{id} [get]
func (ch *CollectionHandler) GetCollectionByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	collection, err := ch.collectionService.GetByID(ctx, collectionUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting collection")
	}

	render.JSON(w, r, ch.dtoBuilder.BuildCollectionDTO(collection))
	return nil
}

// DeleteCollection godoc
// @Summary Delete a collection
// @Description Deletes a collection by its ID. Returns no content on success.
// @Tags collections
// @Security     BearerAuth
// @Param id path string true "Collection ID"
// @Router /collections/{id} [delete]
func (ch *CollectionHandler) DeleteCollection(w http.ResponseWriter, r *http.Request) error {
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	err = ch.collectionService.Delete(r.Context(), collectionUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error deleting collection")
	}

	render.NoContent(w, r)
	return nil
}

// UpdateCollection godoc
// @Summary Update a collection
// @Description Updates a collection by its ID. Returns the updated collection object.
// @Tags collections
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Collection ID"
// @Param title formData string true "Collection title"
// @Param description formData string true "Collection description"
// @Param cover formData file true "Collection cover image"
// @Router /collections/{id} [put]
func (ch *CollectionHandler) UpdateCollection(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "error parsing form")
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	var coverFile multipart.File
	var coverHeader *multipart.FileHeader

	if r.MultipartForm != nil && len(r.MultipartForm.File["cover"]) > 0 {
		coverFile, coverHeader, err = r.FormFile("cover")
		if err != nil {
			return helpers.NewAPIError(http.StatusBadRequest, "error getting cover file")
		}
	}

	preloads := []string{"User"}
	prevCollection, err := ch.collectionService.GetByID(ctx, collectionUUID, preloads...)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting collection")
	}

	_, err = ch.collectionService.UpdateCollection(
		ctx, collectionUUID, services.UpdateCollectionParams{
			UserID:      prevCollection.User.ID,
			Title:       title,
			Description: description,
			CoverHeader: coverHeader,
			Cover:       coverFile,
		},
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error updating collection")
	}

	newCollection, err := ch.collectionService.GetByID(ctx, prevCollection.ID, preloads...)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting collection")
	}

	render.JSON(w, r, ch.dtoBuilder.BuildCollectionDTO(newCollection))
	return nil
}

// GetUsersCollections godoc
// @Summary      Get user's collections
// @Description  Retrieves all collections belonging to the authenticated user
// @Tags         collections
// @Security     BearerAuth
// @Produce      json
// @Router       /collections/users-collections [get]
func (ch *CollectionHandler) GetUsersCollections(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, _ := helpers.GetUserID(ctx)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID format")
	}

	preloads := []string{"Collection"}
	userCollections, err := ch.userCollectionService.Where(
		ctx,
		&models.UserCollection{
			UserID: userUUID,
		},
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting user collections")
	}

	if len(userCollections) == 0 {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, []dto.UserCollectionDTO{})
		return nil
	}

	collections := make([]models.Collection, len(userCollections))
	for i, userCollection := range userCollections {
		collections[i] = userCollection.Collection
	}

	usersCollectionsDTOs := make([]dto.UserCollectionDTO, 0)
	for _, collection := range collections {
		usersCollectionsDTOs = append(usersCollectionsDTOs, ch.dtoBuilder.BuildUserCollectionDTO(&collection))
	}
	render.JSON(w, r, usersCollectionsDTOs)
	return nil
}

// GetCollections godoc
// @Summary Get all collections
// @Description Get all collections. Returns a list of collections.
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param orderBy query string false "Order by (title, created_at)"
// @Param sort query string false "sort order (asc, desc)"
// @Router /collections [get]
func (ch *CollectionHandler) GetCollections(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orderBy := r.URL.Query().Get("orderBy")
	sort := r.URL.Query().Get("sort")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	allowedSortFields := []string{"title", "created_at"}
	defaultField := "created_at"
	preloads := []string{"User"}
	collections, err := ch.collectionService.Where(
		ctx,
		&models.Collection{},
		query.WithPagination(page, limit),
		query.WithSort(orderBy, sort, allowedSortFields, defaultField),
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting collections")
	}

	collectionsDTOs := make([]dto.CollectionDTO, 0)
	for _, collection := range collections {
		collectionsDTOs = append(collectionsDTOs, *ch.dtoBuilder.BuildCollectionDTO(&collection))
	}
	render.JSON(w, r, collectionsDTOs)
	return nil
}

// AddCollectionToUser godoc
// @Summary Add a collection to a user
// @Description Adds a collection to a user. Returns the updated list of collections.
// @Tags collections
// @Security     BearerAuth
// @Produce  json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/add-to-user [post]
func (ch *CollectionHandler) AddCollectionToUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	userID, _ := helpers.GetUserID(ctx)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID format")
	}

	// TODO: change to First with not found error
	userCollectionParams := &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collectionUUID,
	}
	if userCollections, err := ch.userCollectionService.Where(ctx, userCollectionParams); err != nil || len(userCollections) > 0 {
		return helpers.NewAPIError(http.StatusBadRequest, "could not add collection")
	}

	err = ch.userCollectionService.Create(ctx, &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collectionUUID,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "error adding song to collection")
	}

	preloads := []string{"Collection"}
	userCollections, err := ch.userCollectionService.Where(
		ctx,
		&models.UserCollection{
			UserID: userUUID,
		},
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting user collections")
	}

	if len(userCollections) == 0 {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, []dto.UserCollectionDTO{})
		return nil
	}

	collections := make([]models.Collection, len(userCollections))
	for i, userCollection := range userCollections {
		collections[i] = userCollection.Collection
	}

	var usersCollectionsDTOs []dto.UserCollectionDTO
	for _, collection := range collections {
		usersCollectionsDTOs = append(usersCollectionsDTOs, ch.dtoBuilder.BuildUserCollectionDTO(&collection))
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, usersCollectionsDTOs)
	return nil
}

// RemoveCollectionFromUser godoc
// @Summary Remove a collection from a user
// @Description Removes a collection from a user. Returns no content on success.
// @Tags collections
// @Security     BearerAuth
// @Produce  json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/remove-from-user [delete]
func (ch *CollectionHandler) RemoveCollectionFromUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID format")
	}

	userCollection, err := ch.userCollectionService.First(ctx, &models.UserCollection{
		UserID:       userUUID,
		CollectionID: collectionUUID,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error getting user collection")
	}

	err = ch.userCollectionService.Delete(ctx, userCollection.ID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "error removing song from collection")
	}

	render.NoContent(w, r)
	return nil
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
func (ch *CollectionHandler) GetCollectionSongs(w http.ResponseWriter, r *http.Request) error {
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
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

	collectionSongs, err := ch.collectionService.GetCollectionSongs(
		r.Context(),
		collectionUUID,
		search,
		sortBy,
		order,
		page,
		limit,
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "could not retrieve collection songs")
	}
	songDTOs := make([]dto.SongPreviewDTO, 0)
	for _, song := range collectionSongs {
		songDTOs = append(songDTOs, *ch.dtoBuilder.BuildSongPreviewDTO(&song))
	}
	render.JSON(w, r, collectionSongs)
	return nil
}
