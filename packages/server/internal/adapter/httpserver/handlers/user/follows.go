package user

import (
	"errors"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// FollowUser	godoc
// @Summary		Follow user using id param
// @Description	Follow user using id param
// @Tags		user
// @Security	BearerAuth
// @Produce		json
// @Param		id path string true "User to follow ID"
// @Router		/users/{id}/follow [post]
func (uh *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	targetUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid targed user id")
	}
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.BadRequest("invalid targed user id")
	}

	userFollower := &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	}
	if _, err := uh.userFollowerService.First(ctx, userFollower); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return helpers.BadRequest("you are already followed to this user")
		}

		return helpers.InternalServerError("could not follow you to this user")
	}

	if err := uh.userFollowerService.Create(ctx, userFollower); err != nil {
		return helpers.InternalServerError("could not follow this user")
	}

	preloads := []string{"User", "User.Followers", "Follower", "Follower.Followers"}
	userFollowerWithPreloads, _ := uh.userFollowerService.First(ctx, userFollower, preloads...)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, uh.dtoBuilder.BuildUserFollowerDTO(userFollowerWithPreloads))
	return nil
}

// UnfollowUser	godoc
// @Summary		Unfollow user using id param
// @Description	Unfollow user using id param
// @Tags		user
// @Security	BearerAuth
// @Produce		json
// @Param		id path string true "User to unfollow ID"
// @Router		/users/{id}/unfollow [delete]
func (uh *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	targetUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid targed user id")
	}
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.BadRequest("invalid targed user id")
	}

	userFollower, err := uh.userFollowerService.First(ctx, &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	})
	if err != nil {
		return helpers.InternalServerError("trouble with finding out if you followed to this user")
	}

	if err = uh.userFollowerService.Delete(ctx, userFollower.ID); err != nil {
		return helpers.InternalServerError("can not unfollow this user")
	}

	render.NoContent(w, r)
	return nil
}

// IsFollowed godoc
// @Tags	user
// @Security BearerAuth
// @Produce json
// @Param id path string true "The user ID you use to check if you are followed"
// @Router /users/{id}/is-followed [get]
func (uh *UserHandler) IsFollowed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.BadRequest("wrong user id")
	}
	targetUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("wrong target user id")
	}

	userFollowersParams := &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	}
	if _, err := uh.userFollowerService.First(ctx, userFollowersParams); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.JSON(w, r, map[string]bool{"isFollowed": false})
			return nil
		}

		return helpers.InternalServerError("troubles to check if you followed")
	}

	render.JSON(w, r, map[string]bool{"isFollowed": true})
	return nil
}
