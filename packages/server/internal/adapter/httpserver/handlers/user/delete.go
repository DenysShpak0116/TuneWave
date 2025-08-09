package user

import (
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// Delete godoc
// @Summary      Delete user
// @Description  Deletes a user by their ID
// @Tags         user
// @Security     BearerAuth
// @Param        id   path      string true "User ID (UUID format)"
// @Router       /users/{id} [delete]
func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	userUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid user id")
	}

	if err = uh.userService.Delete(r.Context(), userUUID); err != nil {
		return helpers.InternalServerError("failed to delete user")
	}

	render.NoContent(w, r)
	return nil
}
