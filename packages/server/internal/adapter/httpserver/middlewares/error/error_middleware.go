package error

import (
	"log"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/render"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if err := helpers.GetAPIError(r.Context()); err != nil {
			switch e := err.(type) {
			case *helpers.APIError:
				log.Printf("[API ERROR] Status: %d, Message: %s", e.Status, e.Message)
				render.Status(r, e.Status)
				render.JSON(w, r, map[string]any{"error": e.Message})
			default:
				log.Printf("Unhandled error: %v", err)
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]any{"error": "internal server error"})
			}
		}
	})
}
