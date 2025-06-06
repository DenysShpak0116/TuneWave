package result

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ResultHandler struct {
	ResultService         services.ResultService
	CollectionSongService services.CollectionSongService
}

func NewResultHandler(resultService services.ResultService, collectionSongService services.CollectionSongService) *ResultHandler {
	return &ResultHandler{
		ResultService:         resultService,
		CollectionSongService: collectionSongService,
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
func (h *ResultHandler) SendResult(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	var request dto.SendResultRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	results, err := h.ResultService.ProcessUserResults(r.Context(), userUUID, collectionUUID, request)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to process results", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, results)
}

// @Summary Delete user results
// @Description Delete user results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/delete-user-results [delete]
func (h *ResultHandler) DeleteUserResults(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	userResults, err := h.ResultService.Where(context.Background(), &models.Result{
		UserID:           userUUID,
		CollectionSongID: collectionUUID,
	})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get user results", err)
		return
	}
	if len(userResults) == 0 {
		handlers.RespondWithError(w, r, http.StatusNotFound, "No results found for user", nil)
		return
	}

	for _, result := range userResults {
		if err := h.ResultService.Delete(context.Background(), result.ID); err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to delete result", err)
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Results deleted successfully"})
}

// @Summary Get user results
// @Description Get user results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/get-user-results [get]
func (h *ResultHandler) GetUserResults(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	results, err := h.ResultService.GetUserResults(r.Context(), userUUID, collectionUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get results", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, results)
}

// @Summary Get collective results
// @Description Get collective results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/get-results [get]
func (h *ResultHandler) GetCollectiveResults(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	result, err := h.ResultService.GetCollectiveResults(r.Context(), collectionUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get collective results", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, result)
}
