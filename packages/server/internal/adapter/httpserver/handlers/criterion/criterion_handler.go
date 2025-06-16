package criterion

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CriterionHandler struct {
	CriterionService services.CriterionService
}

func NewCriterionHandler(criterionService services.CriterionService) *CriterionHandler {
	return &CriterionHandler{
		CriterionService: criterionService,
	}
}

// r.Get("/", criterionHandler.GetCriterions)
// r.Put("/{id}", criterionHandler.UpdateCriterion)
// r.Delete("/{id}", criterionHandler.DeleteCriterion)

type CreateCriterionRequest struct {
	Name string `json:"name"`
}

// CreateCriterion godoc
// @Summary      Create a new criterion
// @Description  Create a new criterion
// @Tags         criterions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        criterion  body      CreateCriterionRequest  true  "Create criterion"
// @Router       /criterions/ [post]
func (h *CriterionHandler) CreateCriterion(w http.ResponseWriter, r *http.Request) error {
	var request CreateCriterionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	criterion := &models.Criterion{
		Name: request.Name,
	}
	err := h.CriterionService.Create(context.Background(), criterion)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to create criterion")
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, criterion)
	return nil
}

// GetCriterions godoc
// @Summary      Get all criterions
// @Description  Get all criterions
// @Tags         criterions
// @Security     BearerAuth
// @Produce      json
// @Router       /criterions/ [get]
func (h *CriterionHandler) GetCriterions(w http.ResponseWriter, r *http.Request) error {
	criterions, err := h.CriterionService.Where(context.Background(), &models.Criterion{})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get criterions")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, criterions)
	return nil
}

type UpdateCriterionRequest struct {
	Name string `json:"name"`
}

// UpdateCriterion godoc
// @Summary      Update a criterion
// @Description  Update a criterion
// @Tags         criterions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path      string  true  "Criterion ID"
// @Param        criterion  body      UpdateCriterionRequest  true  "Update criterion"
// @Router       /criterions/{id} [put]
func (h *CriterionHandler) UpdateCriterion(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid UUID")
	}

	var request UpdateCriterionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	criterion := &models.Criterion{
		BaseModel: models.BaseModel{
			ID: uuid,
		},
		Name: request.Name,
	}

	newCriterion, err := h.CriterionService.Update(context.Background(), criterion)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to update criterion")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, newCriterion)
	return nil
}

// DeleteCriterion godoc
// @Summary      Delete a criterion
// @Description  Delete a criterion
// @Tags         criterions
// @Security     BearerAuth
// @Param        id  path      string  true  "Criterion ID"
// @Router       /criterions/{id} [delete]
func (h *CriterionHandler) DeleteCriterion(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid UUID")
	}

	err = h.CriterionService.Delete(context.Background(), uuid)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to delete criterion")
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
	return nil
}
