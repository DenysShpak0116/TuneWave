package auth

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// Refresh godoc
// @Summary      Refresh access and refresh tokens
// @Description  Takes a valid refresh token and returns a new pair of access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "Refresh Token"
// @Security     BearerAuth
// @Router       /auth/refresh [post]
func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}

	userID, err := ah.ParseToken(req.RefreshToken)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	accessToken, refreshToken, err := ah.GenerateTokens(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to generate tokens", err)
		return
	}

	user, err := ah.UserService.GetByID(r.Context(), userUUID)

	userDTO := &dtos.UserDTO{}
	if err := dtoMapper.Map(userDTO, user); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to map user", err)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         userDTO,
	})
}
