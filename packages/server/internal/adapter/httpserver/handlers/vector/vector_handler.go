package vector

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/query"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type VectorHandler struct {
	vectorService         services.VectorService
	collectionSongService services.CollectionSongService
	criterionService      services.CriterionService
	dtoBuilder            *dto.DTOBuilder
}

func NewVectorHandler(
	vectorService services.VectorService,
	collectionSongService services.CollectionSongService,
	criterionService services.CriterionService,
	dtoBuilder *dto.DTOBuilder,
) *VectorHandler {
	return &VectorHandler{
		vectorService:         vectorService,
		collectionSongService: collectionSongService,
		criterionService:      criterionService,
		dtoBuilder:            dtoBuilder,
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
func (vh *VectorHandler) GetSongVectors(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection id")
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song id")
	}

	preloads := []string{"Vectors", "Vectors.Criterion"}
	collectionSongs, err := vh.collectionSongService.Where(
		ctx,
		&models.CollectionSong{
			CollectionID: collectionUUID,
			SongID:       songUUID,
		},
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get collection song")
	}

	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return nil
	}

	vectorDTOs := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSongs[0].Vectors {
		vectorDTOs = append(vectorDTOs, vh.dtoBuilder.BuildVectorDTO(vector))
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, vectorDTOs)
	return nil
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
func (vh *VectorHandler) CreateSongVectors(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection id")
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song id")
	}

	var request CreateSongVectorsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	collectionSongs, err := vh.collectionSongService.Where(ctx, &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get collection song")
	}

	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return nil
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
		if err := vh.vectorService.Create(ctx, vector); err != nil {
			return helpers.NewAPIError(http.StatusInternalServerError, "failed to create vector")
		}
	}

	preloads := []string{"Vectors", "Vectors.Criterion"}
	collectionSong, err = vh.collectionSongService.GetByID(ctx, collectionSong.ID, preloads...)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get collection song")
	}

	vectorsDTO := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSong.Vectors {
		vectorsDTO = append(vectorsDTO, vh.dtoBuilder.BuildVectorDTO(vector))
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, vectorsDTO)
	return nil
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
func (vh *VectorHandler) UpdateSongVectors(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection id")
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song id")
	}

	var request UpdateSongVectorsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	for _, vector := range request.Vectors {
		if _, err := vh.vectorService.Update(ctx, &models.Vector{
			BaseModel: models.BaseModel{
				ID: vector.ID,
			},
			Mark:        vector.Mark,
			CriterionID: vector.CriterionID,
		}); err != nil {
			return helpers.NewAPIError(http.StatusInternalServerError, "failed to update vector")
		}
	}

	preloads := []string{"Vectors", "Vectors.Criterion"}
	collectionSongs, err := vh.collectionSongService.Where(
		ctx,
		&models.CollectionSong{
			CollectionID: collectionUUID,
			SongID:       songUUID,
		},
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get collection song")
	}
	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return nil
	}

	vectorsDTO := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSongs[0].Vectors {
		vectorsDTO = append(vectorsDTO, vh.dtoBuilder.BuildVectorDTO(vector))
	}

	render.JSON(w, r, vectorsDTO)
	return nil
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
func (vh *VectorHandler) DeleteSongVectors(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection id")
	}

	songID := chi.URLParam(r, "song-id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song id")
	}

	preloads := []string{
		"Vectors",
	}
	collectionSongs, err := vh.collectionSongService.Where(
		ctx,
		&models.CollectionSong{
			CollectionID: collectionUUID,
			SongID:       songUUID,
		},
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get collection song")
	}

	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{})
		return nil
	}

	for _, vector := range collectionSongs[0].Vectors {
		if err := vh.vectorService.Delete(ctx, vector.ID); err != nil {
			return helpers.NewAPIError(http.StatusInternalServerError, "failed to delete vector")
		}
	}

	render.NoContent(w, r)
	return nil
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
func (vh *VectorHandler) HasAllVectors(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid collection id")
	}

	criterionCount, err := vh.criterionService.CountWhere(ctx, &models.Criterion{})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to count criteria")
	}
	if criterionCount == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, "there is no criteria")
		return nil
	}

	preloads := []string{"Vectors"}
	collectionSongs, err := vh.collectionSongService.Where(
		ctx,
		&models.CollectionSong{CollectionID: collectionUUID},
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get collection song")
	}
	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, "there is no collection song")
		return nil
	}

	for _, collectionSong := range collectionSongs {
		if len(collectionSong.Vectors) < int(criterionCount) {
			render.Status(r, http.StatusOK)
			render.JSON(w, r, map[string]bool{"hasAllVectors": false})
			return nil
		}
	}
	hasAllVectors := true

	render.JSON(w, r, map[string]bool{"hasAllVectors": hasAllVectors})
	return nil
}
