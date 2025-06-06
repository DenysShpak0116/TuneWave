package song

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
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
// @Param sortBy query string false "Sort by field (created_at, title, artist, genre)" default(created_at)
// @Param order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Produce json
// @Router /songs [get]
func (sh *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
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
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get songs", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, songs)
}

// GetByID godoc
// @Summary Get a song by ID
// @Description Get a song by ID
// @Security BearerAuth
// @Tags songs
// @Param id path string true "Song ID"
// @Produce json
// @Router /songs/{id} [get]
func (sh *SongHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
		return
	}

	songDTO, err := sh.SongService.GetFullDTOByID(r.Context(), songUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get song", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, songDTO)
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
// @Failure 401 {object} helpers.ErrorResponse
// @Router /songs [post]
func (sh *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Error parsing form", err)
		return
	}

	userID := r.FormValue("userId")
	userIDuuid, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	title := r.FormValue("title")
	genre := r.FormValue("genre")
	artists := r.Form["artists"]
	tags := r.Form["tags"]

	songFile, songHeader, err := r.FormFile("song")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Song file is required", err)
		return
	}
	defer songFile.Close()

	coverFile, coverHeader, err := r.FormFile("cover")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Cover image is required", err)
		return
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
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to save song", err)
		return
	}

	songDTO, err := sh.SongService.GetFullDTOByID(r.Context(), song.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get song", err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, songDTO)
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
func (sh *SongHandler) Update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Error parsing form", err)
		return
	}

	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
		return
	}

	title := r.FormValue("title")
	genre := r.FormValue("genre")
	artists := r.Form["artists"]
	tags := r.Form["tags"]

	var songFile multipart.File
	var songHeader *multipart.FileHeader
	songFile, songHeader, err = r.FormFile("song")
	if err != nil && err != http.ErrMissingFile {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song file", err)
		return
	}
	if songFile != nil {
		defer songFile.Close()
	}

	var coverFile multipart.File
	var coverHeader *multipart.FileHeader
	coverFile, coverHeader, err = r.FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid cover image", err)
		return
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
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to update song", err)
		return
	}

	updatedSong, err := sh.SongService.GetFullDTOByID(r.Context(), songUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to retrieve updated song", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, updatedSong)
}

// Delete godoc
// @Summary Delete a song
// @Description Delete a song by ID
// @Security BearerAuth
// @Tags songs
// @Param id path string true "Song ID"
// @Produce json
// @Router /songs/{id} [delete]
func (sh *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
		return
	}

	if err := sh.SongService.Delete(context.Background(), songUUID); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to delete song", err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
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
func (sh *SongHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres := sh.SongService.GetGenres(context.Background())
	if len(genres) == 0 {
		render.JSON(w, r, []string{})
		return
	}

	fmt.Printf("genres: %+v\n\n", genres)

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
}
