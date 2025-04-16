package handlers

import (
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/go-chi/render"
)

func RespondWithError(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	render.JSON(w, r, helpers.NewErrorResponse(status, message, errMsg))
}
