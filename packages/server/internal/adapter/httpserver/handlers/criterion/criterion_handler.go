package criterion

import (
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
	criterionService services.CriterionService
}

func NewCriterionHandler(criterionService services.CriterionService) *CriterionHandler {
	return &CriterionHandler{
		criterionService: criterionService,
	}
}

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

	criterion := &models.Criterion{Name: request.Name}
	err := h.criterionService.Create(r.Context(), criterion)
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
	criterions, err := h.criterionService.Where(r.Context(), &models.Criterion{})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get criterions")
	}

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
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
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
	if err := h.criterionService.Update(r.Context(), criterion); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to update criterion")
	}

	render.JSON(w, r, criterion)
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
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid UUID")
	}

	err = h.criterionService.Delete(r.Context(), uuid)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to delete criterion")
	}

	render.NoContent(w, r)
	return nil
}
