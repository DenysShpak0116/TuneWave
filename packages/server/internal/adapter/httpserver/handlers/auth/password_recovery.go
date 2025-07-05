package auth

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/render"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ForgotPassword godoc
// @Summary      Initiate password reset process
// @Description  Sends a password reset link to the user's email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ForgotPasswordRequest true "Email address for password reset"
// @Router       /auth/forgot-password [post]
func (ah *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.BadRequest("invalid request")
	}

	token, err := ah.authService.HandleForgotPassword(req.Email)
	if err != nil {
		return helpers.InternalServerError("failed to send email")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"token": token})
	return nil
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets the user's password using the token received via email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "New password and token"
// @Router       /auth/reset-password [post]
func (ah *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.BadRequest("invalid request")
	}

	if err := ah.authService.HandleResetPassword(req.Token, req.NewPassword); err != nil {
		return helpers.InternalServerError("failed to reset password")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Password reset successfully"})
	return nil
}
