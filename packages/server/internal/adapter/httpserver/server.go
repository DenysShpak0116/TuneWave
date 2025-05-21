package httpserver

import (
	"fmt"
	"log/slog"

	_ "github.com/DenysShpak0116/TuneWave/packages/server/docs"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/chat"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/collection"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/comment"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/criterion"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/result"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/song"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/user"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/vector"
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
	chatHandler *chat.ChatHandler,
	criterionHandler *criterion.CriterionHandler,
	vectorHandler *vector.VectorHandler,
	resultHandler *result.ResultHandler,
) *chi.Mux {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		ExposedHeaders:   []string{"Content-Type", "Content-Length"},
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
		httpSwagger.DocExpansion("false"),
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
		r.Put("/avatar/", userHandler.UpdateAvatar)
		r.Delete("/{id}", userHandler.Delete)
	})

	router.Route("/songs", func(r chi.Router) {
		r.Get("/", songHandler.GetSongs)
		r.Get("/{id}", songHandler.GetByID)
		r.Get("/{id}/is-reacted/{userId}", songHandler.CheckReaction)
		r.Post("/{id}/listen/{userId}", songHandler.ListenSong)

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
		r.Get("/", collectionHandler.GetCollections)
		r.Get("/{id}", collectionHandler.GetCollectionByID)
		r.Get("/{id}/{song-id}/vectors", vectorHandler.GetSongVectors)
		r.Post("/{id}/{song-id}/vectors", vectorHandler.CreateSongVectors)
		r.Put("/{id}/{song-id}/vectors", vectorHandler.UpdateSongVectors)
		r.Delete("/{id}/{song-id}/vectors", vectorHandler.DeleteSongVectors)
		r.Get("/{id}/has-all-vectors", vectorHandler.HasAllVectors)

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

			protected.Post("/{id}/send-results", resultHandler.SendResult)
			protected.Get("/{id}/get-user-results", resultHandler.GetUserResults)
			protected.Get("/{id}/get-results", resultHandler.GetCollectiveResults)
			protected.Delete("/{id}/delete-user-results", resultHandler.DeleteUserResults)

			protected.Post("/", collectionHandler.CreateCollection)
			protected.Put("/{id}", collectionHandler.UpdateCollection)
			protected.Delete("/{id}", collectionHandler.DeleteCollection)
			protected.Get("/users-collections", collectionHandler.GetUsersCollections)
			protected.Post("/{id}/add-to-user", collectionHandler.AddCollectionToUser)
			protected.Delete("/{id}/remove-from-user", collectionHandler.RemoveCollectionFromUser)
		})
	})

	router.Route("/criterions", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

		r.Post("/", criterionHandler.CreateCriterion)
		r.Get("/", criterionHandler.GetCriterions)
		r.Put("/{id}", criterionHandler.UpdateCriterion)
		r.Delete("/{id}", criterionHandler.DeleteCriterion)
	})

	router.Get("/ws/chat", chatHandler.ServeWs)

	router.Route("/chats", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))
		r.Get("/", userHandler.GetChats)
	})

	router.Get("/genres", songHandler.GetGenres)

	return router
}
