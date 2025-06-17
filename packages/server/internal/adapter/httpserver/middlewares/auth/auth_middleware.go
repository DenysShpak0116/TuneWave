package authmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				ctx := helpers.SetAPIError(r.Context(), helpers.NewAPIError(http.StatusUnauthorized, "Authorization header missing or invalid"))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return jwtSecret, nil
			})
			if err != nil {
				ctx := helpers.SetAPIError(r.Context(), helpers.NewAPIError(http.StatusUnauthorized, "Invalid token"))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if exp, ok := claims["exp"].(float64); ok && time.Now().After(time.Unix(int64(exp), 0)) {
					ctx := helpers.SetAPIError(r.Context(), helpers.NewAPIError(http.StatusUnauthorized, "Token expired"))
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				userID, ok := claims["userId"].(string)
				if !ok || userID == "" {
					ctx := helpers.SetAPIError(r.Context(), helpers.NewAPIError(http.StatusUnauthorized, "userId missing in token"))
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				ctx := context.WithValue(r.Context(), "userID", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := helpers.SetAPIError(r.Context(), helpers.NewAPIError(http.StatusUnauthorized, "Invalid token claims"))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
