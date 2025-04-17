package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// Update godoc
// @Summary      Update user
// @Description  Updates a user's profile data (username and profile info) by ID
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string                true  "User ID (UUID format)"
// @Param        user body      dto.UserUpdateRequest true  "Updated user data"
// @Router       /users/{id} [put]
func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "User ID is required"})
		return
	}

	var userUpdateRequest dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	uuidID, err := uuid.Parse(id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid User ID format"})
		return
	}
	userUpdate := &models.User{
		BaseModel: models.BaseModel{
			ID: uuidID,
		},
		Username:    userUpdateRequest.Username,
		ProfileInfo: userUpdateRequest.ProfileInfo,
	}
	updatedUser, err := uh.UserService.Update(context.TODO(), userUpdate)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update user"})
		return
	}

	userDTO := &dtos.UserDTO{}

	if err := dtoMapper.Map(userDTO, updatedUser); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": fmt.Sprintf("Failed to map user to DTO: %v", err)})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"user": userDTO,
	})
}
