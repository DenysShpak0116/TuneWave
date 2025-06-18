package handlers

import (
	"log"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/render"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHandler(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {

			switch e := err.(type) {
			case *helpers.APIError:
				log.Printf("[API ERROR] Status: %d, Message: %s", e.Status, e.Message)
				render.Status(r, e.Status)
				render.JSON(w, r, map[string]any{"error": e.Message})
			default:
				log.Printf("[ERROR] %v", err)
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]any{"error": "internal server error"})
			}
		}
	}
}

