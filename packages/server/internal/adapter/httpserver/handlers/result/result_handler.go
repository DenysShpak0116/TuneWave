package result

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ResultHandler struct {
	resultService         services.ResultService
	collectionSongService services.CollectionSongService
}

func NewResultHandler(resultService services.ResultService, collectionSongService services.CollectionSongService) *ResultHandler {
	return &ResultHandler{
		resultService:         resultService,
		collectionSongService: collectionSongService,
	}
}

// @Summary Send result
// @Description Send result
// @Security BearerAuth
// @Tags result
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Param request body dto.SendResultRequest true "Send result"
// @Router /collections/{id}/send-results [post]
func (vh *ResultHandler) SendResult(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, _ := helpers.GetUserID(ctx)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	var request dto.SendResultRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request body")
	}

	results, err := vh.resultService.ProcessUserResults(ctx, userUUID, collectionUUID, request)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to process results")
	}

	render.JSON(w, r, results)
	return nil
}

// @Summary Delete user results
// @Description Delete user results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/delete-user-results [delete]
func (vh *ResultHandler) DeleteUserResults(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, _ := helpers.GetUserID(ctx)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	userResults, err := vh.resultService.Where(ctx, &models.Result{
		UserID:           userUUID,
		CollectionSongID: collectionUUID,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get user results")
	}
	if len(userResults) == 0 {
		return helpers.NewAPIError(http.StatusNotFound, "no results found for user")
	}

	ids := make([]uuid.UUID, 0)
	for _, result := range userResults {
		ids = append(ids, result.ID)
	}
	if err := vh.resultService.Delete(ctx, ids...); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to delete result")
	}

	render.JSON(w, r, map[string]string{"message": "Results deleted successfully"})
	return nil
}

// @Summary Get user results
// @Description Get user results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/get-user-results [get]
func (vh *ResultHandler) GetUserResults(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, _ := helpers.GetUserID(ctx)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	results, err := vh.resultService.GetUserResults(ctx, userUUID, collectionUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get results")
	}

	render.JSON(w, r, results)
	return nil
}

// @Summary Get collective results
// @Description Get collective results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/get-results [get]
func (vh *ResultHandler) GetCollectiveResults(w http.ResponseWriter, r *http.Request) error {
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection ID")
	}

	result, err := vh.resultService.GetCollectiveResults(r.Context(), collectionUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get collective results")
	}

	render.JSON(w, r, result)
	return nil
}
