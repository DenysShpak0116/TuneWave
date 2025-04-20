package user

import (
	"context"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GetByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Router /users/{id} [get]
func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userIDuuid, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid user ID", err)
		return
	}
	user, err := uh.UserService.GetFullDTOByID(context.Background(), userIDuuid)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusNotFound, "user not found", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}
