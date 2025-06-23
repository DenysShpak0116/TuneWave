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
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid collection id")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "song-id"))
	if err != nil {
		return helpers.BadRequest("invalid song id")
	}

	collectionSongParams := &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	}
	preloads := []string{"Vectors", "Vectors.Criterion"}
	collectionSong, err := vh.collectionSongService.First(ctx, collectionSongParams, preloads...)
	if err != nil {
		return helpers.InternalServerError("failed to get collection song")
	}

	vectorDTOs := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSong.Vectors {
		vectorDTOs = append(vectorDTOs, vh.dtoBuilder.BuildVectorDTO(vector))
	}
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
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid collection id")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "song-id"))
	if err != nil {
		return helpers.BadRequest("invalid song id")
	}

	var request CreateSongVectorsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.BadRequest("invalid request")
	}

	collectionSong, err := vh.collectionSongService.First(ctx, &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	})
	if err != nil {
		return helpers.InternalServerError("failed to get collection song")
	}

	vectors := make([]*models.Vector, len(request.Vectors))
	for i, vector := range request.Vectors {
		vectors[i] = &models.Vector{
			Mark:             vector.Mark,
			CriterionID:      vector.CriterionID,
			CollectionSongID: collectionSong.ID,
		}
	}

	if err := vh.vectorService.Create(ctx, vectors...); err != nil {
		return helpers.InternalServerError("failed to create vector")
	}

	preloads := []string{"Vectors", "Vectors.Criterion"}
	collectionSong, err = vh.collectionSongService.GetByID(ctx, collectionSong.ID, preloads...)
	if err != nil {
		return helpers.InternalServerError("failed to get collection song")
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
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid collection id")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "song-id"))
	if err != nil {
		return helpers.BadRequest("invalid song id")
	}

	var request UpdateSongVectorsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.BadRequest("invalid request")
	}

	for _, vector := range request.Vectors {
		if err := vh.vectorService.Update(ctx, &models.Vector{
			BaseModel: models.BaseModel{
				ID: vector.ID,
			},
			Mark:        vector.Mark,
			CriterionID: vector.CriterionID,
		}); err != nil {
			return helpers.InternalServerError("failed to update vector")
		}
	}

	collectionSongParams := &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	}
	preloads := []string{"Vectors", "Vectors.Criterion"}
	collectionSong, err := vh.collectionSongService.First(ctx, collectionSongParams, preloads...)
	if err != nil {
		return helpers.InternalServerError("failed to get collection song")
	}

	vectorsDTO := make([]dto.VectorDTO, 0)
	for _, vector := range collectionSong.Vectors {
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
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid collection id")
	}
	songUUID, err := uuid.Parse(chi.URLParam(r, "song-id"))
	if err != nil {
		return helpers.BadRequest("invalid song id")
	}

	collectionSongParams := &models.CollectionSong{
		CollectionID: collectionUUID,
		SongID:       songUUID,
	}
	preloads := []string{"Vectors"}
	collectionSong, err := vh.collectionSongService.First(ctx, collectionSongParams, preloads...)
	if err != nil {
		return helpers.InternalServerError("failed to get collection song")
	}

	ids := make([]uuid.UUID, 0)
	for _, vector := range collectionSong.Vectors {
		ids = append(ids, vector.ID)
	}
	if err := vh.vectorService.Delete(ctx, ids...); err != nil {
		return helpers.InternalServerError("failed to delete vector")
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
	collectionUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid collection id")
	}

	criterionCount, err := vh.criterionService.CountWhere(ctx, &models.Criterion{})
	if err != nil {
		return helpers.InternalServerError("failed to count criteria")
	}
	if criterionCount == 0 {
		return helpers.NotFound("there is no criteria")
	}

	preloads := []string{"Vectors"}
	collectionSongs, err := vh.collectionSongService.Where(
		ctx,
		&models.CollectionSong{CollectionID: collectionUUID},
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.InternalServerError("failed to get collection song")
	}
	if len(collectionSongs) == 0 {
		return helpers.NotFound("there is no collection song")
	}

	for _, collectionSong := range collectionSongs {
		if len(collectionSong.Vectors) < int(criterionCount) {
			render.JSON(w, r, map[string]bool{"hasAllVectors": false})
			return nil
		}
	}

	render.JSON(w, r, map[string]bool{"hasAllVectors": true})
	return nil
}
