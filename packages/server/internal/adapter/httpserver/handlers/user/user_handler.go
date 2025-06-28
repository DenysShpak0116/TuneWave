package user

import (
	"net/http"
	"strconv"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/query"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/render"
)

type UserHandler struct {
	userService         services.UserService
	userFollowerService services.UserFollowerService
	userReactionService services.UserReactionService
	messageService      services.MessageService
	dtoBuilder          *dto.DTOBuilder
}

func NewUserHandler(
	userService services.UserService,
	userFollowerService services.UserFollowerService,
	userReactionService services.UserReactionService,
	messageService services.MessageService,
	dtoBuilder *dto.DTOBuilder,
) *UserHandler {
	return &UserHandler{
		userService:         userService,
		userFollowerService: userFollowerService,
		userReactionService: userReactionService,
		messageService:      messageService,
		dtoBuilder:          dtoBuilder,
	}
}

// GetAll godoc
// @Summary      Get all users
// @Description  Get all users with pagination
// @Security     BearerAuth
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        page  query  int  false  "Page number"  default(1)
// @Param        limit query  int  false  "Number of users per page"  default(10)
// @Router       /users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	users, err := h.userService.Where(ctx, &models.User{}, query.WithPagination(page, limit))
	if err != nil {
		return helpers.InternalServerError("failed to get users")
	}

	render.JSON(w, r, users)
	return nil
}
