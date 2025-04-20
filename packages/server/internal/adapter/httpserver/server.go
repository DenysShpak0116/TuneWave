package httpserver

import (
	"fmt"
	"log/slog"

	_ "github.com/DenysShpak0116/TuneWave/packages/server/docs"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/collection"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/comment"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/song"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/user"
	authmiddleware "github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/middlewares/auth"
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
	collectionHandler *collection.CollectionHandler,
) *chi.Mux {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
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

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/logout", authHandler.Logout)
		r.Post("/forgot-password", authHandler.ForgotPassword)
		r.Post("/reset-password", authHandler.ResetPassword)
		r.Post("/refresh", authHandler.Refresh)

		r.Get("/google/callback", authHandler.GoogleCallback)
		r.Get("/google", authHandler.GoogleAuth)
	})

	router.Route("/users", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

		r.Get("/", userHandler.GetAll)
		r.Get("/{id}", userHandler.GetByID)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
	})

	router.Route("/songs", func(r chi.Router) {
		r.Get("/", songHandler.GetSongs)
		r.Get("/{id}", songHandler.GetByID)

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

			protected.Post("/", songHandler.Create)
			protected.Put("/{id}", songHandler.Update)
			protected.Delete("/{id}", songHandler.Delete)

			protected.Post("/{id}/reaction", songHandler.SetReaction)

			protected.Post("/{id}/add-to-collection", songHandler.AddToCollection)
		})
	})

	router.Route("/comments", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

		r.Post("/", commentHandler.CreateComment)
		r.Delete("/{id}", commentHandler.DeleteComment)
	})

	router.Route("/collections", func(r chi.Router) {
		r.Get("/{id}", collectionHandler.GetCollectionByID)

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

			protected.Post("/", collectionHandler.CreateCollection)
			protected.Put("/{id}", collectionHandler.UpdateCollection)
			protected.Delete("/{id}", collectionHandler.DeleteCollection)
		})
	})

	return router
}
