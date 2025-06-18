package authmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]any{"error": "Authorization header missing or invalid"})
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
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]any{"error": "Invalid token"})
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if exp, ok := claims["exp"].(float64); ok {
					if time.Now().After(time.Unix(int64(exp), 0)) {
						render.Status(r, http.StatusUnauthorized)
						render.JSON(w, r, map[string]any{"error": "Token expired"})
						return
					}
				}

				userID, ok := claims["userId"].(string)
				if !ok || userID == "" {
					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, map[string]any{"error": "userId missing in token"})
					return
				}

				ctx := context.WithValue(r.Context(), "userID", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]any{"error": "Invalid token claims"})
		})
	}
}
