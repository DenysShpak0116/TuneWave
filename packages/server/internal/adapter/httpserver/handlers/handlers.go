package handlers

import (
	"net/http"

	"log"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/render"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHandler(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		switch e := err.(type) {
		case *helpers.APIError:
			log.Printf("[API ERROR] Status: %d, Message: %s", e.Status, e.Message)
			render.Status(r, e.Status)
			render.JSON(w, r, map[string]any{"error": e.Message})

		default:
			log.Printf("Error: %v", e)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]any{"error": "internal server error"})
		}
	}
}
