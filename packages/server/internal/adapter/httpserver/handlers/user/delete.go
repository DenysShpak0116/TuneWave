package user

import (
	"context"
	"net/http"

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
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "User ID is required"})
		return nil
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid User ID format"})
		return nil
	}

	err = uh.userService.Delete(context.TODO(), uuidID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete user"})
		return nil
	}

	render.Status(r, http.StatusNoContent)
	return nil
}
