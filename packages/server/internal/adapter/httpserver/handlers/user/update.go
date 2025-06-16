package user

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
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
func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "User ID is required"})
		return nil
	}

	var userUpdateRequest dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return nil
	}
	defer r.Body.Close()

	uuidID, err := uuid.Parse(id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid User ID format"})
		return nil
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
		return nil
	}

	updatedUser, err = uh.UserService.GetByID(context.Background(), updatedUser.ID, "Followers")

	dtoBuilder := dto.NewDTOBuilder().
		SetCountUserFollowersFunc(func(userID uuid.UUID) int64 {
			return uh.UserService.GetUserFollowersCount(context.Background(), userID)
		})
	userDTO := dtoBuilder.BuildUserDTO(updatedUser)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"user": userDTO,
	})
	return nil
}

// UpdateAvatar godoc
// @Summary      Update user avatar
// @Description  Updates a user's avatar
// @Tags         user
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Avatar file"
// @Router 	 /users/avatar/ [put]
func (uh *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid form data")
	}

	var pfpFile multipart.File
	var pfpHeader *multipart.FileHeader
	pfpFile, pfpHeader, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid cover image")
	}
	if pfpFile != nil {
		defer pfpFile.Close()
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		return helpers.NewAPIError(http.StatusBadRequest, "user id is required")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid User ID format")
	}

	err = uh.UserService.UpdateUserPfp(context.TODO(), services.UpdatePfpParams{
		UserID:    userUUID,
		Pfp:       pfpFile,
		PfpHeader: pfpHeader,
	})
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update avatar"})
		return nil
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Avatar updated successfully"})
	return nil
}
