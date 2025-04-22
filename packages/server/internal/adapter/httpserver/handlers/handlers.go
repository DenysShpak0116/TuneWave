package handlers

import (
	"net/http"

	"log"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/render"
)

func RespondWithError(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	log.Printf("Error: %v, Status: %d, Message: %s", err, status, message)

	var ErrMsg string
	if err != nil {
		ErrMsg = err.Error()
	} else {
		ErrMsg = message
	}

	render.Status(r, status)
	render.JSON(w, r, helpers.NewErrorResponse(status, message, ErrMsg))
}
