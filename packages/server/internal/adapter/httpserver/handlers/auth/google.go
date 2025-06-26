package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/markbates/goth/gothic"
)

// GoogleAuth godoc
// @Summary Start Google authentication
// @Description Redirects to Google OAuth 2.0 login
// @Tags auth
// @Accept  json
// @Produce  json
// @Router /auth/google [get]
func (ah *AuthHandler) GoogleAuth(res http.ResponseWriter, req *http.Request) error {
	gothic.BeginAuthHandler(res, req)
	return nil
}

type UserWithNickname struct {
	Nicknames []struct {
		Value string `json:"value"`
	} `json:"nicknames"`
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handles the callback after Google authentication, fetches user info, and generates access and refresh tokens.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param code query string true "Google OAuth code"
// @Router /auth/google/callback [get]
func (ah *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return helpers.BadRequest("failed to get user data")
	}

	var currentUser *models.User
	fetchedUser, err := ah.userService.First(ctx, &models.User{Email: user.Email})
	if errors.Is(err, service.ErrNotFound) {
		nickname, err := fetchGoogleNickname(user.AccessToken)
		if err != nil {
			return helpers.InternalServerError("failed to fetch nickname")
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
		if err := ah.userService.Create(ctx, currentUser); err != nil {
			return helpers.InternalServerError("failed to create user")
		}
	} else if err != nil {
		return helpers.InternalServerError("failed to get user data")
	} else {
		currentUser = fetchedUser
	}

	accessToken, refreshToken, err := ah.GenerateTokens(currentUser.ID.String())
	if err != nil {
		return helpers.InternalServerError("failed to generate tokens")
	}

	authData := map[string]any{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         ah.dtoBuilder.BuildUserDTO(currentUser),
	}
	authJSON, err := json.Marshal(authData)
	if err != nil {
		return helpers.InternalServerError("failed to encode auth data")
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

	http.Redirect(w, r, "http://localhost:5173/", http.StatusSeeOther)
	return nil
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

	var userModel UserWithNickname
	if err := json.NewDecoder(resp.Body).Decode(&userModel); err != nil {
		return "", err
	}

	if userModel.Nicknames == nil || len(userModel.Nicknames) == 0 {
		return "", err
	}
	return userModel.Nicknames[0].Value, nil
}
