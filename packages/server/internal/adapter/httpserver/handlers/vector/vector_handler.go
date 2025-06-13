package vector

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

type VectorHandler struct {
	VactorService         services.VectorService
	CollectionSongService services.CollectionSongService
	CriterionService      services.CriterionService
}

func NewVectorHandler(
	vectorService services.VectorService,
	collectionSongService services.CollectionSongService,
	criterionService services.CriterionService,
) *VectorHandler {
	return &VectorHandler{
		VactorService:         vectorService,
		CollectionSongService: collectionSongService,
		CriterionService:      criterionService,
	}
}

// GetSongVectors godoc
// @Summary Get song vectors
// @Description Get song vectors
// @Security BearerAuth
// @Tags vectors
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Param song-id path string true "Song ID"
// @Router /collections/{id}/{song-id}/vectors [get]
func (h *VectorHandler) GetSongVectors(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid collection id", err)
		return
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid song id", err)
		return
	}

	collectionSongs, err := h.CollectionSongService.Where(r.Context(), &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	}, "Vectors", "Vectors.Criterion")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to get collection song", err)
		return
	}

	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return
	}
	collectionSong := collectionSongs[0]

	dtoBuilder := dto.NewDTOBuilder()

	vectorsDTO := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSong.Vectors {
		vectorsDTO = append(vectorsDTO, dtoBuilder.BuildVectorDTO(vector))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, vectorsDTO)
}

type CreateSongVectorsRequest struct {
	Vectors []struct {
		CriterionID uuid.UUID `json:"criterionId"`
		Mark        string    `json:"mark"`
	} `json:"vectors"`
}

// CreateSongVectors godoc
// @Summary Create song vectors
// @Description Create song vectors
// @Security BearerAuth
// @Tags vectors
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Param song-id path string true "Song ID"
// @Param vectors body CreateSongVectorsRequest true "Vectors"
// @Router /collections/{id}/{song-id}/vectors [post]
func (h *VectorHandler) CreateSongVectors(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid collection id", err)
		return
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid song id", err)
		return
	}

	var request CreateSongVectorsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid request", err)
		return
	}

	collectionSongs, err := h.CollectionSongService.Where(r.Context(), &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to get collection song", err)
		return
	}

	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return
	}
	collectionSong := &collectionSongs[0]

	vectors := make([]*models.Vector, len(request.Vectors))
	for i, vector := range request.Vectors {
		vectors[i] = &models.Vector{
			Mark:             vector.Mark,
			CriterionID:      vector.CriterionID,
			CollectionSongID: collectionSong.ID,
		}
	}

	for _, vector := range vectors {
		if err := h.VactorService.Create(r.Context(), vector); err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to create vector", err)
			return
		}
	}

	collectionSong, err = h.CollectionSongService.GetByID(r.Context(), collectionSong.ID, "Vectors", "Vectors.Criterion")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to get collection song", err)
		return
	}

	dtoBuilder := dto.NewDTOBuilder()

	vectorsDTO := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSong.Vectors {
		vectorsDTO = append(vectorsDTO, dtoBuilder.BuildVectorDTO(vector))
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, vectorsDTO)
}

type UpdateSongVectorsRequest struct {
	Vectors []struct {
		ID          uuid.UUID `json:"id"`
		CriterionID uuid.UUID `json:"criterionId"`
		Mark        string    `json:"mark"`
	} `json:"vectors"`
}

// UpdateSongVectors godoc
// @Summary Update song vectors
// @Description Update song vectors
// @Security BearerAuth
// @Tags vectors
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Param song-id path string true "Song ID"
// @Param vectors body UpdateSongVectorsRequest true "Vectors"
// @Router /collections/{id}/{song-id}/vectors [put]
func (h *VectorHandler) UpdateSongVectors(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid collection id", err)
		return
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid song id", err)
		return
	}

	var request UpdateSongVectorsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid request", err)
		return
	}

	for _, vector := range request.Vectors {
		if _, err := h.VactorService.Update(r.Context(), &models.Vector{
			BaseModel: models.BaseModel{
				ID: vector.ID,
			},
			Mark:        vector.Mark,
			CriterionID: vector.CriterionID,
		}); err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to update vector", err)
			return
		}
	}

	collectionSongs, err := h.CollectionSongService.Where(r.Context(), &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	}, "Vectors", "Vectors.Criterion")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to get collection song", err)
		return
	}
	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return
	}
	collectionSong := collectionSongs[0]

	dtoBuilder := dto.NewDTOBuilder()

	vectorsDTO := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSong.Vectors {
		vectorsDTO = append(vectorsDTO, dtoBuilder.BuildVectorDTO(vector))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, vectorsDTO)
}

// DeleteSongVectors godoc
// @Summary Delete song vectors
// @Description Delete song vectors
// @Security BearerAuth
// @Tags vectors
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Param song-id path string true "Song ID"
// @Router /collections/{id}/{song-id}/vectors [delete]
func (h *VectorHandler) DeleteSongVectors(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid collection id", err)
		return
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid song id", err)
		return
	}

	collectionSongs, err := h.CollectionSongService.Where(r.Context(), &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	}, "Vectors")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to get collection song", err)
		return
	}

	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return
	}
	collectionSong := collectionSongs[0]

	for _, vector := range collectionSong.Vectors {
		if err := h.VactorService.Delete(r.Context(), vector.ID); err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to delete vector", err)
			return
		}
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// HasAllVectors godoc
// @Summary Check if all vectors are present
// @Description Check if all vectors are present
// @Security BearerAuth
// @Tags vectors
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/has-all-vectors [get]
func (h *VectorHandler) HasAllVectors(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid collection id", err)
		return
	}

	criterionCount, err := h.CriterionService.CountWhere(context.Background(), &models.Criterion{})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to count criteria", err)
		return
	}
	if criterionCount == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, "there is no criteria")
		return
	}

	collectionSongs, err := h.CollectionSongService.Where(r.Context(), &models.CollectionSong{
		CollectionID: collectionUUID,
	}, "Vectors")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "failed to get collection song", err)
		return
	}
	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, "there is no collection song")
		return
	}

	for _, collectionSong := range collectionSongs {
		if len(collectionSong.Vectors) < int(criterionCount) {
			render.Status(r, http.StatusOK)
			render.JSON(w, r, map[string]bool{"hasAllVectors": false})
			return
		}
	}
	hasAllVectors := true

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]bool{"hasAllVectors": hasAllVectors})
}
