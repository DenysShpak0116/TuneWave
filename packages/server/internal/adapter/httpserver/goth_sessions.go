package httpserver

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func InitGothicSessionStore() {
	key := "something-very-secret"
	maxAge := 86400 * 30
	isProd := false

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "google", nil
	}
}
