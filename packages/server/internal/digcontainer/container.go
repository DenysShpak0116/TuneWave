package digcontainer

import (
	"os"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/comment"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/song"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/user"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/logger/slogpretty"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/repository"
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
	container.Provide(func(db *gorm.DB) *repository.UserRepository {
		return repository.NewUserRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.TokenRepository {
		return repository.NewTokenRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.SongRepository {
		return repository.NewSongRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.AuthorRepository {
		return repository.NewAuthorRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.SongAuthorRepository {
		return repository.NewSongAuthorRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.TagRepository {
		return repository.NewTagRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.SongTagRepository {
		return repository.NewSongTagRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.UserReactionRepository {
		return repository.NewUserReactionRepository(db)
	})

	container.Provide(func(db *gorm.DB) *repository.CommentRepository {
		return repository.NewCommentRepository(db)
	})
	// service
	container.Provide(func(userRepo *repository.UserRepository) *service.UserService {
		return service.NewUserService(userRepo)
	})

	container.Provide(func(cfg *config.Config) *service.MailService {
		return service.NewMailService(
			cfg.Mail.StmpServer,
			cfg.Mail.SmtpPort,
			cfg.Mail.FromMail,
			cfg.Mail.FromPassword,
		)
	})

	container.Provide(func(
		mailService *service.MailService,
		tokenRepo *repository.TokenRepository,
		userService *service.UserService,
	) *service.AuthService {
		return service.NewAuthService(mailService, tokenRepo, userService)
	})

	container.Provide(func(
		songRepo *repository.SongRepository,
		fileStorage *repository.FileStorage,
		authorRepository *repository.AuthorRepository,
		songAuthorRepository *repository.SongAuthorRepository,
		tagRepository *repository.TagRepository,
		songTagRepository *repository.SongTagRepository,
		userReactionRepository *repository.UserReactionRepository,
	) *songservice.SongService {
		return songservice.NewSongService(
			songRepo,
			fileStorage,
			authorRepository,
			songAuthorRepository,
			tagRepository,
			songTagRepository,
			userReactionRepository,
		)
	})

	container.Provide(func(
		commentRepo *repository.CommentRepository,
	) *service.CommentService {
		return service.NewCommentService(commentRepo)
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
		cfg *config.Config,
		log *slog.Logger,
		authHandler *auth.AuthHandler,
		userHandler *user.UserHandler,
		songHandler *song.SongHandler,
		commentHandler *comment.CommentHandler,
	) *chi.Mux {
		return httpserver.NewRouter(log, cfg, authHandler, userHandler, songHandler, commentHandler)
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
