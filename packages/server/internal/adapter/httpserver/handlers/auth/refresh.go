package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// Refresh godoc
// @Summary      Refresh access and refresh tokens
// @Description  Takes a valid refresh token and returns a new pair of access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Router       /auth/refresh [post]
func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	authCookie, err := r.Cookie("authData")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "authData cookie not found", err)
		return
	}

	authDataBytes, err := base64.URLEncoding.DecodeString(authCookie.Value)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Failed to decode authData", err)
		return
	}

	var authData map[string]interface{}
	if err := json.Unmarshal(authDataBytes, &authData); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Failed to parse authData", err)
		return
	}

	refreshToken, ok := authData["refreshToken"].(string)
	if !ok || refreshToken == "" {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "refreshToken not found in authData", nil)
		return
	}

	userID, err := helpers.ParseToken(ah.JWTSecret, refreshToken)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	accessToken, newRefreshToken, err := ah.GenerateTokens(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to generate tokens", err)
		return
	}

	user, err := ah.UserService.GetByID(ctx, userUUID, "Followers")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	newAuthData := map[string]interface{}{
		"refreshToken": newRefreshToken,
	}
	authJSON, err := json.Marshal(newAuthData)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to encode auth data", err)
		return
	}

	authBase64 := base64.URLEncoding.EncodeToString(authJSON)

	http.SetCookie(w, &http.Cookie{
		Name:     "authData",
		Value:    authBase64,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	dtoBuilder := dto.NewDTOBuilder().
		SetCountUserFollowersFunc(func(userID uuid.UUID) int64 {
			return ah.UserService.GetUserFollowersCount(ctx, userID)
		})

	render.JSON(w, r, map[string]interface{}{
		"accessToken": accessToken,
		"user":        dtoBuilder.BuildUserDTO(user),
	})
}
