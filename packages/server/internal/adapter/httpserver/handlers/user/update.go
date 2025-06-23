package user

import (
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
	ctx := r.Context()
	userUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}

	var userUpdateRequest dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateRequest); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request payload")
	}
	defer r.Body.Close()

	updatedUser := &models.User{
		BaseModel: models.BaseModel{
			ID: userUUID,
		},
		Username:    userUpdateRequest.Username,
		ProfileInfo: userUpdateRequest.ProfileInfo,
	}
	if err := uh.userService.Update(ctx, updatedUser); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "Failed to update user")
	}

	render.JSON(w, r, uh.dtoBuilder.BuildUserDTO(updatedUser))
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
	ctx := r.Context()
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

	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid User ID format")
	}

	err = uh.userService.UpdateUserPfp(ctx, services.UpdatePfpParams{
		UserID:    userUUID,
		Pfp:       pfpFile,
		PfpHeader: pfpHeader,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "Failed to update avatar")
	}

	render.JSON(w, r, map[string]string{"message": "Avatar updated successfully"})
	return nil
}
