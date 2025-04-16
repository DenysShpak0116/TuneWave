package auth

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/go-chi/render"
)

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}

	userID, err := ah.ParseToken(req.RefreshToken)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	accessToken, refreshToken, err := ah.GenerateTokens(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to generate tokens", err)
		return
	}

	render.JSON(w, r, map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
