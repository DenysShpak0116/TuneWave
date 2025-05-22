package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with email, password, and username. Returns the created user object.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body dto.RegisterRequest true "User registration data"
// @Router /auth/register [post]
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	existingUsers, err := ah.UserService.Where(ctx, &models.User{Email: req.Email})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to check existing users", err)
		return
	}

	if len(existingUsers) > 0 {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "User already exists", nil)
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user := &models.User{
		Email:           req.Email,
		Username:        req.Username,
		Role:            "user",
		IsGoogleAccount: false,
		PasswordHash:    hash,
		ProfilePicture:  "https://photosrush.com/wp-content/uploads/dark-aesthetic-anime-pfp-girl-1.jpg",
	}

	if err := ah.UserService.Create(ctx, user); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	userDTO, err := ah.UserService.GetFullDTOByID(ctx, user.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get user DTO", err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, userDTO)
}

// Login godoc
// @Summary Login an existing user
// @Description Logs in an existing user with email and password, and returns access and refresh tokens.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param login body dto.LoginRequest true "User login data"
// @Router /auth/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}

	ctx := context.Background()

	users, err := ah.UserService.Where(ctx, &models.User{Email: req.Email})
	if err != nil || len(users) == 0 {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid credentials", nil)
		return
	}

	user := users[0]

	if user.IsGoogleAccount {
		handlers.RespondWithError(w, r, http.StatusForbidden, "This email is associated with a Google account. Please log in with Google.", nil)
		return
	}

	fmt.Println("Password Hash:", user.PasswordHash)

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		handlers.RespondWithError(w, r, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	accessToken, refreshToken, err := ah.GenerateTokens(user.ID.String())
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to generate tokens", err)
		return
	}

	userData, err := ah.UserService.GetByID(ctx, user.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get user DTO", err)
		return
	}

	userDTO := &dtos.UserDTO{
		ID:             userData.ID,
		Username:       userData.Username,
		Role:           userData.Role,
		ProfilePicture: userData.ProfilePicture,
		ProfileInfo:    userData.ProfileInfo,
	}

	authData := map[string]any{
		"refreshToken": refreshToken,
	}

	authJSON, err := json.Marshal(authData)
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
	})

	render.JSON(w, r, map[string]interface{}{
		"accessToken": accessToken,
		"user":        userDTO,
	})
}

// Logout godoc
// @Summary Logout a user
// @Description Logs out the user by invalidating their authentication token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Router /auth/logout [post]
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		handlers.RespondWithError(w, r, http.StatusUnauthorized, "Missing Authorization header", nil)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := helpers.ParseToken(ah.JWTSecret, tokenStr)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}

	uuidParsed, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}
	users, err := ah.UserService.Where(ctx, &models.User{BaseModel: models.BaseModel{ID: uuidParsed}})
	if err != nil || len(users) == 0 {
		handlers.RespondWithError(w, r, http.StatusUnauthorized, "User not found", err)
		return
	}

	user := users[0]

	if user.IsGoogleAccount {
		if err := gothic.Logout(w, r); err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to logout Google user", err)
			return
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
}
