package digcontainer

import (
	"os"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/chat"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/collection"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/comment"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/song"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/user"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/ws"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/logger/slogpretty"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/repository"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/songservice"

	"log/slog"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	httpserver.InitGothicSessionStore()

	container.Provide(config.MustLoad)
	container.Provide(setupLogger)
	container.Provide(repository.NewGORMDB)

	// repository
	container.Provide(repository.NewFileStorage)
	container.Provide(repository.NewRepository[models.User])
	container.Provide(repository.NewRepository[models.Token])
	container.Provide(repository.NewRepository[models.Song])
	container.Provide(repository.NewRepository[models.Author])
	container.Provide(repository.NewRepository[models.SongAuthor])
	container.Provide(repository.NewRepository[models.Tag])
	container.Provide(repository.NewRepository[models.SongTag])
	container.Provide(repository.NewRepository[models.UserReaction])
	container.Provide(repository.NewRepository[models.Comment])
	container.Provide(repository.NewRepository[models.Collection])
	container.Provide(repository.NewRepository[models.CollectionSong])
	container.Provide(repository.NewRepository[models.Chat])
	container.Provide(repository.NewRepository[models.Message])

	// service
	container.Provide(service.NewMailService)
	container.Provide(songservice.NewSongService)
	container.Provide(service.NewUserService)
	container.Provide(service.NewAuthService)
	container.Provide(service.NewCommentService)
	container.Provide(service.NewCollectionService)
	container.Provide(service.NewChatService)
	container.Provide(service.NewMessageService)
	container.Provide(ws.NewHubManager)

	// handlers
	container.Provide(chat.NewChatHandler)
	container.Provide(auth.NewAuthHandler)
	container.Provide(user.NewUserHandler)
	container.Provide(song.NewSongHandler)
	container.Provide(comment.NewCommentHandler)
	container.Provide(collection.NewCollectionHandler)

	container.Provide(httpserver.NewRouter)

	return container
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	switch cfg.Env {
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
