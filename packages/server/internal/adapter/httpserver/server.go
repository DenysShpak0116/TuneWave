package httpserver

import (
	"fmt"
	"log/slog"

	_ "github.com/DenysShpak0116/TuneWave/packages/server/docs"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/comment"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/song"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/user"
	mwLogger "github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/middlewares/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	log *slog.Logger,
	cfg *config.Config,
	authHandler *auth.AuthHandler,
	userHandler *user.UserHandler,
	songHandler *song.SongHandler,
	commentHandler *comment.CommentHandler,
) *chi.Mux {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	router := chi.NewRouter()

	router.Use(corsHandler.Handler)
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.Http.Port)),
	))

	router.Post("/auth/login", authHandler.Login)
	router.Post("/auth/register", authHandler.Register)
	router.Post("/auth/logout", authHandler.Logout)
	router.Post("/auth/forgot-password", authHandler.ForgotPassword)
	router.Post("/auth/reset-password", authHandler.ResetPassword)
	router.Get("/auth/refresh", authHandler.Refresh)

	router.Get("/auth/google/callback", authHandler.GoogleCallback)
	router.Get("/auth/google", authHandler.GoogleAuth)

	router.Route("/users", func(r chi.Router) {
		r.Get("/{id}", userHandler.GetByID)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
	})

	router.Route("/songs", func(r chi.Router) {
		r.Get("/", songHandler.GetSongs)
		r.Get("/{id}", songHandler.GetByID)
		r.Post("/", songHandler.Create)
		r.Put("/{id}", songHandler.Update)
		r.Delete("/{id}", songHandler.Delete)

		r.Post("/{id}/reaction", songHandler.SetReaction)
	})

	router.Route("/comments", func(r chi.Router) {
		r.Post("/", commentHandler.CreateComment)
		r.Delete("/{id}", commentHandler.DeleteComment)
	})

	return router
}
