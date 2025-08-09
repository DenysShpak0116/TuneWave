package helpers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MuxSetURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
