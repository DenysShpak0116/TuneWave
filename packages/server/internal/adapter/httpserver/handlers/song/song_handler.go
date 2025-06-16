package song

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SongHandler struct {
	SongService           services.SongService
	CollectionSongService services.CollectionSongService
}

func NewSongHandler(songService services.SongService, collectionSongService services.CollectionSongService) *SongHandler {
	return &SongHandler{
		SongService:           songService,
		CollectionSongService: collectionSongService,
	}
}

// GetSongs godoc
// @Summary Get all songs
// @Description Get all songs
// @Tags songs
// @Param search query string false "Search by title, artist, or genre"
// @Param sortBy query string false "Sort by field (created_at, title, artist, genre, listenings)" default(created_at)
// @Param order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Produce json
// @Router /songs [get]
func (sh *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) error {
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

	songs, err := sh.SongService.GetSongs(context.Background(), search, sortBy, order, page, limit)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get songs")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, songs)
	return nil
}

// GetByID godoc
// @Summary Get a song by ID
// @Description Get a song by ID
// @Security BearerAuth
// @Tags songs
// @Param id path string true "Song ID"
// @Produce json
// @Router /songs/{id} [get]
func (sh *SongHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	songDTO, err := sh.SongService.GetByID(r.Context(), songUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get song")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, songDTO)
	return nil
}

// Create godoc
// @Summary Create a new song
// @Description Create a new song
// @Security BearerAuth
// @Tags songs
// @Accept multipart/form-data
// @Produce json
// @Param userId formData string true "User ID"
// @Param title formData string true "Song title"
// @Param genre formData string true "Song genre"
// @Param artists formData []string true "Artists" collectionFormat(multi)
// @Param tags formData []string true "Tags" collectionFormat(multi)
// @Param song formData file true "Song file"
// @Param cover formData file true "Cover image"
// @Router /songs [post]
func (sh *SongHandler) Create(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "error parsing form")
	}

	userID := r.FormValue("userId")
	userIDuuid, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}

	title := r.FormValue("title")
	genre := r.FormValue("genre")
	artists := r.Form["artists"]
	tags := r.Form["tags"]

	songFile, songHeader, err := r.FormFile("song")
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "song file is required")
	}
	defer songFile.Close()

	coverFile, coverHeader, err := r.FormFile("cover")
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "cover image is required")
	}
	defer coverFile.Close()

	song, err := sh.SongService.SaveSong(context.Background(), services.SaveSongParams{
		UserID:      userIDuuid,
		Title:       title,
		Genre:       genre,
		Artists:     artists,
		Tags:        tags,
		Song:        songFile,
		SongHeader:  songHeader,
		Cover:       coverFile,
		CoverHeader: coverHeader,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to save song")
	}

	songDTO, err := sh.SongService.GetByID(r.Context(), song.ID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get song")
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, songDTO)
	return nil
}

// Update godoc
// @Summary Update a song
// @Description Update a song by ID
// @Security BearerAuth
// @Tags songs
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Song ID"
// @Param title formData string false "Song title"
// @Param genre formData string false "Song genre"
// @Param artists formData []string false "Artists" collectionFormat(multi)
// @Param tags formData []string false "Tags" collectionFormat(multi)
// @Param song formData file false "Updated song file"
// @Param cover formData file false "Updated cover image"
// @Router /songs/{id} [put]
func (sh *SongHandler) Update(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "error parsing form")
	}

	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	title := r.FormValue("title")
	genre := r.FormValue("genre")
	artists := r.Form["artists"]
	tags := r.Form["tags"]

	var songFile multipart.File
	var songHeader *multipart.FileHeader
	songFile, songHeader, err = r.FormFile("song")
	if err != nil && err != http.ErrMissingFile {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song file")
	}
	if songFile != nil {
		defer songFile.Close()
	}

	var coverFile multipart.File
	var coverHeader *multipart.FileHeader
	coverFile, coverHeader, err = r.FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid cover image")
	}
	if coverFile != nil {
		defer coverFile.Close()
	}

	err = sh.SongService.UpdateSong(context.Background(), services.UpdateSongParams{
		SongID:      songUUID,
		Title:       title,
		Genre:       genre,
		Artists:     artists,
		Tags:        tags,
		Song:        songFile,
		SongHeader:  songHeader,
		Cover:       coverFile,
		CoverHeader: coverHeader,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to update song")
	}

	updatedSong, err := sh.SongService.GetByID(r.Context(), songUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "Failed to retrieve updated song")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, updatedSong)
	return nil
}

// Delete godoc
// @Summary Delete a song
// @Description Delete a song by ID
// @Security BearerAuth
// @Tags songs
// @Param id path string true "Song ID"
// @Produce json
// @Router /songs/{id} [delete]
func (sh *SongHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	if err := sh.SongService.Delete(context.Background(), songUUID); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to delete song")
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
	return nil
}

type genrePreview struct {
	GenreName  string `json:"genreName"`
	GenreCover string `json:"genreCover"`
}

// GetGenres godoc
// @Summary Get genres previews
// @Description Get genres previews
// @Tags songs
// @Produce json
// @Router /genres [get]
func (sh *SongHandler) GetGenres(w http.ResponseWriter, r *http.Request) error {
	genres := sh.SongService.GetGenres(context.Background())
	if len(genres) == 0 {
		render.JSON(w, r, []string{})
		return nil
	}

	genrePreviews := make([]genrePreview, 0)
	for _, genre := range genres {
		song, err := sh.SongService.GetGenresMostPopularSong(context.Background(), genre)
		fmt.Printf("song: %+v\n\n", song)
		if err != nil || song == nil {
			genrePreviews = append(genrePreviews, genrePreview{
				GenreName:  genre,
				GenreCover: "",
			})
			continue
		}
		genrePreviews = append(genrePreviews, genrePreview{
			GenreName:  genre,
			GenreCover: song.CoverURL,
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, genrePreviews)
	return nil
}
