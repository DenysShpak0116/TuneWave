package user

import (
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
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
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}

	// TODO: change to First logic
	userFollowers, _ := uh.userFollowerService.Where(ctx, &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	})
	if len(userFollowers) > 0 {
		return helpers.NewAPIError(http.StatusBadRequest, "you are already followed to this user")
	}

	userFollower := &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	}
	if err = uh.userFollowerService.Create(ctx, userFollower); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "could not follow this user")
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
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}

	userFollower, err := uh.userFollowerService.First(ctx, &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "trouble with finding out if you followed to this user")
	}

	if err = uh.userFollowerService.Delete(ctx, userFollower.ID); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "can not unfollow this user")
	}

	render.NoContent(w, r)
	return nil
}

// IsFollowed godoc
// @Tags		user
// @Security	BearerAuth
// @Produce json
// @Param id path string true "The user ID you use to check if you are followed"
// @Router /users/{id}/is-followed [get]
func (uh *UserHandler) IsFollowed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "wrong user id")
	}
	targetUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "wrong target user id")
	}

	// TODO: change to First logic
	if userFollower, err := uh.userFollowerService.Where(ctx, &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	}); err != nil || len(userFollower) < 1 {
		if err != nil {
			return helpers.NewAPIError(http.StatusInternalServerError, "troubles to check if you followed")
		}

		render.JSON(w, r, map[string]bool{"isFollowed": false})
		return nil
	}

	render.JSON(w, r, map[string]bool{"isFollowed": true})
	return nil
}
