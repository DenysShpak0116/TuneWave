package digcontainer

import (
	"os"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/collection"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/comment"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/song"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/user"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/logger/slogpretty"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/repository"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/songservice"
	"github.com/go-chi/chi/v5"

	"log/slog"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	httpserver.InitGothicSessionStore()

	container.Provide(config.MustLoad)
	container.Provide(func(cfg *config.Config) *slog.Logger {
		return setupLogger(cfg.Env)
	})

	container.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		return repository.NewGORMDB(cfg.StoragePath)
	})

	// repository
	container.Provide(func(cfg *config.Config) *repository.FileStorage {
		return repository.NewFileStorage(cfg.AWS.Region, cfg.AWS.AccessKey, cfg.AWS.SecretKey, cfg.AWS.Bucket)
	})
	container.Provide(func(
		db *gorm.DB,
	) (port.Repository[models.User],
		port.Repository[models.Token],
		port.Repository[models.Song],
		port.Repository[models.Author],
		port.Repository[models.SongAuthor],
		port.Repository[models.Tag],
		port.Repository[models.SongTag],
		port.Repository[models.UserReaction],
		port.Repository[models.Comment],
		port.Repository[models.Collection],
		port.Repository[models.CollectionSong]) {
		return repository.NewRepository[models.User](db),
			repository.NewRepository[models.Token](db),
			repository.NewRepository[models.Song](db),
			repository.NewRepository[models.Author](db),
			repository.NewRepository[models.SongAuthor](db),
			repository.NewRepository[models.Tag](db),
			repository.NewRepository[models.SongTag](db),
			repository.NewRepository[models.UserReaction](db),
			repository.NewRepository[models.Comment](db),
			repository.NewRepository[models.Collection](db),
			repository.NewRepository[models.CollectionSong](db)
	})

	// service

	container.Provide(func(cfg *config.Config) *service.MailService {
		return service.NewMailService(
			cfg.Mail.StmpServer,
			cfg.Mail.SmtpPort,
			cfg.Mail.FromMail,
			cfg.Mail.FromPassword,
		)
	})

	container.Provide(func(
		songRepo port.Repository[models.Song],
		fileStorage *repository.FileStorage,
		authorRepository port.Repository[models.Author],
		songAuthorRepository port.Repository[models.SongAuthor],
		tagRepository port.Repository[models.Tag],
		songTagRepository port.Repository[models.SongTag],
		userReactionRepository port.Repository[models.UserReaction],
		collectionSongRepository port.Repository[models.CollectionSong],
	) *songservice.SongService {
		return songservice.NewSongService(
			songRepo,
			fileStorage,
			authorRepository,
			songAuthorRepository,
			tagRepository,
			songTagRepository,
			userReactionRepository,
			collectionSongRepository,
		)
	})

	container.Provide(func(userRepo port.Repository[models.User], songService *songservice.SongService) *service.UserService {
		return service.NewUserService(userRepo, songService)
	})

	container.Provide(func(
		mailService *service.MailService,
		tokenRepo port.Repository[models.Token],
		userService *service.UserService,
	) *service.AuthService {
		return service.NewAuthService(mailService, tokenRepo, userService)
	})

	container.Provide(func(
		commentRepo port.Repository[models.Comment],
	) *service.CommentService {
		return service.NewCommentService(commentRepo)
	})

	container.Provide(func(
		collectionRepo port.Repository[models.Collection],
		fileStorage *repository.FileStorage,
		collectionSongRepository port.Repository[models.CollectionSong],
	) *service.CollectionService {
		return service.NewCollectionService(collectionRepo, fileStorage, collectionSongRepository)
	})
	// handlers
	container.Provide(func(
		authService *service.AuthService,
		userService *service.UserService,
		cfg *config.Config,
	) *auth.AuthHandler {
		return auth.NewAuthHandler(
			authService,
			userService,
			cfg.Google.ClientID,
			cfg.Google.ClientSecret,
			cfg.JwtSecret,
		)
	})
	container.Provide(func(
		userService *service.UserService,
	) *user.UserHandler {
		return user.NewUserHandler(userService)
	})
	container.Provide(func(
		songService *songservice.SongService,
	) *song.SongHandler {
		return song.NewSongHandler(songService)
	})

	container.Provide(func(
		commentService *service.CommentService,
	) *comment.CommentHandler {
		return comment.NewCommentHandler(commentService)
	})

	container.Provide(func(
		collectionService *service.CollectionService,
	) *collection.CollectionHandler {
		return collection.NewCollectionHandler(collectionService)
	})

	container.Provide(func(
		cfg *config.Config,
		log *slog.Logger,
		authHandler *auth.AuthHandler,
		userHandler *user.UserHandler,
		songHandler *song.SongHandler,
		commentHandler *comment.CommentHandler,
		collectionHandler *collection.CollectionHandler,
	) *chi.Mux {
		return httpserver.NewRouter(
			log,
			cfg,
			authHandler,
			userHandler,
			songHandler,
			commentHandler,
			collectionHandler,
		)
	})

	return container
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
