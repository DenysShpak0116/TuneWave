package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/go-chi/render"
)

// ForgotPassword godoc
// @Summary      Initiate password reset process
// @Description  Sends a password reset link to the user's email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ForgotPasswordRequest true "Email address for password reset"
// @Router       /auth/forgot-password [post]
func (ah *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, map[string]string{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	token, err := ah.AuthService.HandleForgotPassword(req)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": fmt.Sprintf("failed to send email: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"token": token})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets the user's password using the token received via email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "New password and token"
// @Router       /auth/reset-password [post]
func (ah *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, map[string]string{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if err := ah.AuthService.HandleResetPassword(req); err != nil {
		render.JSON(w, r, map[string]string{"error": fmt.Sprintf("failed to reset password: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Password reset successfully"})
}
