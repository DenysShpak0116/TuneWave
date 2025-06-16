package httpserver

import (
	"fmt"
	"log/slog"

	_ "github.com/DenysShpak0116/TuneWave/packages/server/docs"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
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
		r.Get("/", handlers.MakeHandler(userHandler.GetAll))
		r.Get("/{id}/collections", handlers.MakeHandler(userHandler.GetUserCollections))

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

			protected.Get("/{id}", handlers.MakeHandler(userHandler.GetByID))
			protected.Post("/{id}/follow", handlers.MakeHandler(userHandler.FollowUser))
			protected.Delete("/{id}/unfollow", handlers.MakeHandler(userHandler.UnfollowUser))
			protected.Get("/{id}/is-followed", handlers.MakeHandler(userHandler.IsFollowed))
			protected.Put("/{id}", handlers.MakeHandler(userHandler.Update))
			protected.Put("/avatar/", handlers.MakeHandler(userHandler.UpdateAvatar))
			protected.Delete("/{id}", handlers.MakeHandler(userHandler.Delete))
		})
	})

	router.Route("/songs", func(r chi.Router) {
		r.Get("/", handlers.MakeHandler(songHandler.GetSongs))
		r.Get("/{id}", handlers.MakeHandler(songHandler.GetByID))
		r.Get("/{id}/is-reacted/{userId}", handlers.MakeHandler(songHandler.CheckReaction))
		r.Post("/{id}/listen/{userId}", handlers.MakeHandler(songHandler.ListenSong))

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

			protected.Post("/", handlers.MakeHandler(songHandler.Create))
			protected.Put("/{id}", handlers.MakeHandler(songHandler.Update))
			protected.Delete("/{id}", handlers.MakeHandler(songHandler.Delete))

			protected.Post("/{id}/reaction", handlers.MakeHandler(songHandler.SetReaction))

			protected.Post("/{id}/add-to-collection", handlers.MakeHandler(songHandler.AddToCollection))
			protected.Delete("/{id}/remove-from-collection", handlers.MakeHandler(songHandler.RemoveFromCollection))
		})
	})

	router.Route("/comments", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

		r.Post("/", handlers.MakeHandler(commentHandler.CreateComment))
		r.Delete("/{id}", handlers.MakeHandler(commentHandler.DeleteComment))
	})

	router.Route("/collections", func(r chi.Router) {
		r.Get("/", handlers.MakeHandler(collectionHandler.GetCollections))
		r.Get("/{id}", handlers.MakeHandler(collectionHandler.GetCollectionByID))
		r.Get("/{id}/songs", handlers.MakeHandler(collectionHandler.GetCollectionSongs))
		r.Get("/{id}/{song-id}/vectors", handlers.MakeHandler(vectorHandler.GetSongVectors))
		r.Post("/{id}/{song-id}/vectors", handlers.MakeHandler(vectorHandler.CreateSongVectors))
		r.Put("/{id}/{song-id}/vectors", handlers.MakeHandler(vectorHandler.UpdateSongVectors))
		r.Delete("/{id}/{song-id}/vectors", handlers.MakeHandler(vectorHandler.DeleteSongVectors))
		r.Get("/{id}/has-all-vectors", handlers.MakeHandler(vectorHandler.HasAllVectors))

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

			protected.Post("/{id}/send-results", handlers.MakeHandler(resultHandler.SendResult))
			protected.Get("/{id}/get-user-results", handlers.MakeHandler(resultHandler.GetUserResults))
			protected.Get("/{id}/get-results", handlers.MakeHandler(resultHandler.GetCollectiveResults))
			protected.Delete("/{id}/delete-user-results", handlers.MakeHandler(resultHandler.DeleteUserResults))

			protected.Post("/", handlers.MakeHandler(collectionHandler.CreateCollection))
			protected.Put("/{id}", handlers.MakeHandler(collectionHandler.UpdateCollection))
			protected.Delete("/{id}", handlers.MakeHandler(collectionHandler.DeleteCollection))
			protected.Get("/users-collections", handlers.MakeHandler(collectionHandler.GetUsersCollections))
			protected.Post("/{id}/add-to-user", handlers.MakeHandler(collectionHandler.AddCollectionToUser))
			protected.Delete("/{id}/remove-from-user", handlers.MakeHandler(collectionHandler.RemoveCollectionFromUser))
		})
	})

	router.Route("/criterions", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))

		r.Post("/", handlers.MakeHandler(criterionHandler.CreateCriterion))
		r.Get("/", handlers.MakeHandler(criterionHandler.GetCriterions))
		r.Put("/{id}", handlers.MakeHandler(criterionHandler.UpdateCriterion))
		r.Delete("/{id}", handlers.MakeHandler(criterionHandler.DeleteCriterion))
	})

	router.Get("/ws/chat", chatHandler.ServeWs)

	router.Route("/chats", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware([]byte(cfg.JwtSecret)))
		r.Get("/", handlers.MakeHandler(userHandler.GetChats))
	})

	router.Get("/genres", handlers.MakeHandler(songHandler.GetGenres))

	return router
}
