package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/render"
)

type UserHandler struct {
	UserService         services.UserService
	UserFollowerService services.UserFollowerService
	UserReactionService services.UserReactionService
}

func NewUserHandler(
	userService services.UserService,
	userFollowerService services.UserFollowerService,
	userReactionService services.UserReactionService,
) *UserHandler {
	return &UserHandler{
		UserService:         userService,
		UserFollowerService: userFollowerService,
		UserReactionService: userReactionService,
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
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	users, err := h.UserService.GetUsers(context.Background(), page, limit)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
}
