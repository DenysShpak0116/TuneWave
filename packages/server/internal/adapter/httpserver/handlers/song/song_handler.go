package song

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

type SongHandler struct {
	songService           services.SongService
	collectionSongService services.CollectionSongService
	userReactionService   services.UserReactionService
	commentService        services.CommentService
	dtoBuilder            *dto.DTOBuilder
}

func NewSongHandler(
	songService services.SongService,
	collectionSongService services.CollectionSongService,
	userReactionService services.UserReactionService,
	commentService services.CommentService,
	dtoBuilder *dto.DTOBuilder,
) *SongHandler {
	return &SongHandler{
		songService:           songService,
		collectionSongService: collectionSongService,
		userReactionService:   userReactionService,
		commentService:        commentService,
		dtoBuilder:            dtoBuilder,
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
	ctx := r.Context()

	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "created_at"
	}
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "desc"
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	params := services.SearchSongsParams{
		Search: search,
		SortBy: sortBy,
		Order:  order,
		Page:   page,
		Limit:  limit,
	}
	preloads := []string{"Authors", "Authors.Author"}
	songs, err := sh.songService.GetSongs(ctx, params, preloads...)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get songs")
	}

	songDTOs := make([]dto.SongPreviewDTO, 0)
	for _, song := range songs {
		songDTOs = append(songDTOs, *sh.dtoBuilder.BuildSongPreviewDTO(&song))
	}
	render.JSON(w, r, songDTOs)
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
	ctx := r.Context()

	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	preloads := []string{
		"Authors",
		"Authors.Author",
		"SongTags",
		"SongTags.Tag",
		"User",
	}
	song, err := sh.songService.GetByID(ctx, songUUID, preloads...)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get song")
	}

	render.JSON(w, r, sh.dtoBuilder.BuildSongDTO(song))
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
	ctx := r.Context()

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

	song, err := sh.songService.SaveSong(ctx, services.SaveSongParams{
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

	songDTO, err := sh.songService.GetByID(ctx, song.ID)
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
	ctx := r.Context()

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

	err = sh.songService.UpdateSong(ctx, services.UpdateSongParams{
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

	updatedSong, err := sh.songService.GetByID(ctx, songUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "Failed to retrieve updated song")
	}

	render.JSON(w, r, sh.dtoBuilder.BuildSongDTO(updatedSong))
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
	ctx := r.Context()
	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	if err := sh.songService.Delete(ctx, songUUID); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to delete song")
	}

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
	ctx := r.Context()
	genres := sh.songService.GetGenres(ctx)
	if len(genres) == 0 {
		render.JSON(w, r, []string{})
		return nil
	}

	genrePreviews := make([]genrePreview, 0)
	for _, genre := range genres {
		song, err := sh.songService.GetGenresMostPopularSong(ctx, genre)
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

	render.JSON(w, r, genrePreviews)
	return nil
}

// GetSongComments godoc
// @Summary Get song comments
// @Description Get song comments
// @Tags songs
// @Produce json
// @Param id path string true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Router /songs/{id}/comments [get]
func (sh *SongHandler) GetSongComments(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	songID := chi.URLParam(r, "id")
	songUUID, err := uuid.Parse(songID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	preloads := []string{"User"}
	comments, err := sh.commentService.Where(
		ctx,
		&models.Comment{SongID: songUUID},
		query.WithPagination(page, limit),
		query.WithPreloads(preloads...),
	)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "could not retrieve comments")
	}

	commentDTOs := make([]dto.CommentDTO, 0)
	for _, comment := range comments {
		commentDTOs = append(commentDTOs, *sh.dtoBuilder.BuildCommentDTO(&comment))
	}
	render.JSON(w, r, commentDTOs)
	return nil
}
