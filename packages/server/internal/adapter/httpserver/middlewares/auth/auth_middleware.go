package authmiddleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				handlers.RespondWithError(w, r, http.StatusUnauthorized, "Authorization header missing or invalid", fmt.Errorf("authorization header missing or invalid"))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return jwtSecret, nil
			})

			log.Printf("Parsed token: %v", token)

			if err != nil || !token.Valid {
				log.Printf("Invalid token: %v", err)
				handlers.RespondWithError(w, r, http.StatusUnauthorized, "Invalid token", err)
				return
			}

			log.Printf("Going to check claims: %v", token.Claims)
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				log.Printf("Claims: %v", claims)
				if exp, ok := claims["exp"].(float64); ok {
					log.Printf("Token expiration time: %v", exp)
					if time.Now().After(time.Unix(int64(exp), 0)) {
						log.Printf("Token expired")
						handlers.RespondWithError(w, r, http.StatusUnauthorized, "Token expired", fmt.Errorf("token expired"))
						return
					}
				}
				ctx := context.WithValue(r.Context(), "claims", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			handlers.RespondWithError(w, r, http.StatusUnauthorized, "Invalid token claims", fmt.Errorf("invalid token claims"))
		})
	}
}

func GetClaimsFromContext(r *http.Request) jwt.MapClaims {
	if claims, ok := r.Context().Value("claims").(jwt.MapClaims); ok {
		return claims
	}
	return nil
}
