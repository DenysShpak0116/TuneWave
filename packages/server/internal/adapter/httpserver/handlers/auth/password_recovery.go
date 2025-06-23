package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
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
func (ah *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	var req dto.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	token, err := ah.authService.HandleForgotPassword(req)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to send email")
	}

	render.JSON(w, r, map[string]string{"token": token})
	return nil
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets the user's password using the token received via email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "New password and token"
// @Router       /auth/reset-password [post]
func (ah *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	if err := ah.authService.HandleResetPassword(req); err != nil {
		render.JSON(w, r, map[string]string{"error": fmt.Sprintf("failed to reset password: %v", err)})
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to reset password")
	}

	render.JSON(w, r, map[string]string{"message": "Password reset successfully"})
	return nil
}
