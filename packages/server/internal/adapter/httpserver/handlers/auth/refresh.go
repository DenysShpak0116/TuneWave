package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

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
func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authCookie, err := r.Cookie("authData")
	if err != nil {
		return helpers.BadRequest("authData cookie not found")
	}

	authDataBytes, err := base64.URLEncoding.DecodeString(authCookie.Value)
	if err != nil {
		return helpers.BadRequest("failed to decode authData")
	}

	var authData map[string]any
	if err := json.Unmarshal(authDataBytes, &authData); err != nil {
		return helpers.BadRequest("failed to parse authData")
	}

	refreshToken, ok := authData["refreshToken"].(string)
	if !ok || refreshToken == "" {
		return helpers.BadRequest("refreshToken not found in authData")
	}

	userID, err := helpers.ParseToken(ah.jwtSecret, refreshToken)
	if err != nil {
		return helpers.NewAPIError(http.StatusUnauthorized, "invalid or expired refresh token")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.BadRequest("invalid user ID")
	}

	accessToken, newRefreshToken, err := ah.GenerateTokens(userID)
	if err != nil {
		return helpers.InternalServerError("failed to generate tokens")
	}

	preloads := []string{"Followers"}
	user, err := ah.userService.GetByID(ctx, userUUID, preloads...)
	if err != nil {
		return helpers.InternalServerError("failed to get user")
	}

	newAuthData := map[string]any{
		"refreshToken": newRefreshToken,
	}
	authJSON, err := json.Marshal(newAuthData)
	if err != nil {
		return helpers.InternalServerError("Failed to encode auth data")
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

	render.JSON(w, r, map[string]interface{}{
		"accessToken": accessToken,
		"user":        ah.dtoBuilder.BuildUserDTO(user),
	})
	return nil
}
