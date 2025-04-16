package digcontainer

import (
	"os"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/logger/slogpretty"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/repository"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/go-chi/chi"

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
	container.Provide(func(db *gorm.DB) *repository.UserRepository {
		return repository.NewUserRepository(db)
	})

	// service
	container.Provide(func(userRepo *repository.UserRepository) *service.UserService {
		return service.NewUserService(userRepo)
	})

	// handlers
	container.Provide(func(cfg *config.Config, userService *service.UserService) *auth.AuthHandler {
		return auth.NewAuthHandler(userService, cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.JwtSecret)
	})
	container.Provide(func(
		cfg *config.Config,
		log *slog.Logger,
		authHandler *auth.AuthHandler,
	) *chi.Mux {
		return httpserver.NewRouter(log, cfg, authHandler)
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
