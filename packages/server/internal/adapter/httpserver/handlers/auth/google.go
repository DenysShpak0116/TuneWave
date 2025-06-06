package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/markbates/goth/gothic"
)

// GoogleAuth godoc
// @Summary Start Google authentication
// @Description Redirects to Google OAuth 2.0 login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Router /auth/google [get]
func (ah *AuthHandler) GoogleAuth(res http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(res, req)
}

type UserWithNickname struct {
	Nicknames []struct {
		Value string `json:"value"`
	} `json:"nicknames"`
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handles the callback after Google authentication, fetches user info, and generates access and refresh tokens.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param code query string true "Google OAuth code"
// @Router /auth/google/callback [get]
func (ah *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Failed to get user data", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	users, err := ah.UserService.Where(ctx, &models.User{Email: user.Email})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get user data", err)
		return
	}

	var currentUser *models.User
	if len(users) > 0 {
		currentUser = &users[0]
	} else {
		nickname, err := fetchGoogleNickname(user.AccessToken)
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to fetch nickname", err)
			return
		}
		if nickname == "" {
			nickname = user.Name
		}

		currentUser = &models.User{
			Email:           user.Email,
			Username:        nickname,
			Role:            "user",
			IsGoogleAccount: true,
			PasswordHash:    "",
			ProfileInfo:     "",
			ProfilePicture:  user.AvatarURL,
		}

		if err := ah.UserService.Create(ctx, currentUser); err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to create user", err)
			return
		}
	}

	accessToken, refreshToken, err := ah.GenerateTokens(currentUser.ID.String())
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to generate tokens", err)
		return
	}

	userDTO, err := ah.UserService.GetFullDTOByID(ctx, currentUser.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get user DTO", err)
		return
	}

	authData := map[string]any{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         userDTO,
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

	http.Redirect(w, r, "http://localhost:5173/", http.StatusSeeOther)
}

func fetchGoogleNickname(token string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://people.googleapis.com/v1/people/me?personFields=nicknames",
		nil,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response body:", resp.Body)
	var userModel UserWithNickname
	if err := json.NewDecoder(resp.Body).Decode(&userModel); err != nil {
		return "", err
	}

	if userModel.Nicknames == nil || len(userModel.Nicknames) == 0 {
		return "", err
	}
	return userModel.Nicknames[0].Value, nil
}
