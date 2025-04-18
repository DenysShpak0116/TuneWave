package song

import (
	"context"
	"net/http"
	"strconv"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/render"
)

type SongHandler struct {
	SongService services.SongService
}

func NewSongHandler(songService services.SongService) *SongHandler {
	return &SongHandler{
		SongService: songService,
	}
}

func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
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

	songs, err := h.SongService.GetSongs(context.Background(), search, sortBy, order, page, limit)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get songs", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, songs)
}

func (h *SongHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting a song by ID
	w.Write([]byte("Get song by ID"))
}

func (h *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a new song
	w.Write([]byte("Create a new song"))
}

func (h *SongHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating a song
	w.Write([]byte("Update a song"))
}

func (h *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting a song
	w.Write([]byte("Delete a song"))
}
