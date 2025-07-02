package digcontainer

import (
	"os"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/auth"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/chat"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/collection"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/comment"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/criterion"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/result"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/song"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/user"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/vector"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/ws"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/logger/slogpretty"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/repository"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"

	"log/slog"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(config.MustLoad)

	container.Invoke(func(cfg *config.Config) {
		httpserver.InitGothicSessionStore(cfg.Google.GothicSessionKay, cfg.Google.MaxSessionAge, cfg.Env == "prod")
	})

	container.Provide(setupPrettySlog)
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
	container.Provide(repository.NewRepository[models.Criterion])
	container.Provide(repository.NewRepository[models.Vector])
	container.Provide(repository.NewRepository[models.Result])
	container.Provide(repository.NewRepository[models.UserCollection])
	container.Provide(repository.NewRepository[models.UserFollower])

	// service
	container.Provide(service.NewMailService)
	container.Provide(service.NewSongService)
	container.Provide(service.NewUserService)
	container.Provide(service.NewAuthService)
	container.Provide(service.NewCommentService)
	container.Provide(service.NewCollectionService)
	container.Provide(service.NewCollectionSongService)
	container.Provide(service.NewChatService)
	container.Provide(service.NewMessageService)
	container.Provide(service.NewCriterionService)
	container.Provide(service.NewVectorService)
	container.Provide(service.NewResultService)
	container.Provide(service.NewUserCollectionService)
	container.Provide(service.NewUserFollowerService)
	container.Provide(service.NewUserReactionService)
	container.Provide(ws.NewHubManager)

	// handlers
	container.Provide(func(userService services.UserService, userReactionService services.UserReactionService) *dto.DTOBuilder {
		return dto.NewDTOBuilder(userService, userReactionService)
	})

	container.Provide(chat.NewChatHandler)
	container.Provide(auth.NewAuthHandler)
	container.Provide(user.NewUserHandler)
	container.Provide(song.NewSongHandler)
	container.Provide(comment.NewCommentHandler)
	container.Provide(collection.NewCollectionHandler)
	container.Provide(criterion.NewCriterionHandler)
	container.Provide(vector.NewVectorHandler)
	container.Provide(result.NewResultHandler)

	container.Provide(httpserver.NewRouter)

	return container
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
