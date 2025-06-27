package httpserver

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func InitGothicSessionStore(gothicSessionKey string, maxSessionAge int, isProd bool) {
	store := sessions.NewCookieStore([]byte(gothicSessionKey))
	store.MaxAge(maxSessionAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "google", nil
	}
}
