package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

// Register		godoc
// @Summary		Register a new user
// @Description	Registers a new user with email, password, and username. Returns the created user object.
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		user body dto.RegisterRequest true "User registration data"
// @Router		/auth/register [post]
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	existingUsers, err := ah.userService.Where(ctx, &models.User{Email: req.Email})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to check existing users")
	}

	if len(existingUsers) > 0 {
		return helpers.NewAPIError(http.StatusBadRequest, "user already exists")
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to hash password")
	}

	user := &models.User{
		Email:           req.Email,
		Username:        req.Username,
		Role:            "user",
		IsGoogleAccount: false,
		PasswordHash:    hash,
		ProfilePicture:  "https://photosrush.com/wp-content/uploads/dark-aesthetic-anime-pfp-girl-1.jpg",
	}

	if err := ah.userService.Create(ctx, user); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to create user")
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ah.dtoBuilder.BuildUserDTO(user))
	return nil
}

// Login		godoc
// @Summary		Login an existing user
// @Description	Logs in an existing user with email and password, and returns access and refresh tokens.
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		login body dto.LoginRequest true "User login data"
// @Router		/auth/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	user, err := ah.userService.First(ctx, &models.User{Email: req.Email})
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid credentials")
	}

	if user.IsGoogleAccount {
		return helpers.NewAPIError(http.StatusForbidden, "This email is associated with a Google account. Please log in with Google.")
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		return helpers.NewAPIError(http.StatusUnauthorized, "invalid credentials")
	}

	accessToken, refreshToken, err := ah.GenerateTokens(user.ID.String())
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to generate tokens")
	}

	authData := map[string]any{
		"refreshToken": refreshToken,
	}

	authJSON, err := json.Marshal(authData)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to encode auth data")
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

	render.JSON(w, r, map[string]any{
		"accessToken": accessToken,
		"user":        ah.dtoBuilder.BuildUserDTO(user),
	})
	return nil
}

// Logout godoc
// @Summary Logout a user
// @Description Logs out the user by invalidating their authentication token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Router /auth/logout [post]
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return helpers.NewAPIError(http.StatusUnauthorized, "missing Authorization header")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := helpers.ParseToken(ah.jwtSecret, tokenStr)
	if err != nil {
		return helpers.NewAPIError(http.StatusUnauthorized, "invalid or expired token")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID format")
	}
	user, err := ah.userService.First(ctx, &models.User{BaseModel: models.BaseModel{ID: userUUID}})
	if err != nil {
		return helpers.NewAPIError(http.StatusUnauthorized, "user not found")
	}

	if user.IsGoogleAccount {
		if err := gothic.Logout(w, r); err != nil {
			return helpers.NewAPIError(http.StatusInternalServerError, "failed to logout Google user")
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "authData",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	render.JSON(w, r, map[string]string{
		"message": "Successfully logged out",
	})
	return nil
}
