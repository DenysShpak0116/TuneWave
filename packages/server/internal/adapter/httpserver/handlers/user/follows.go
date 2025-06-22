package user

import (
	"context"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/query"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// TODO: figure out  response body
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
	targetID := chi.URLParam(r, "id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}

	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid auth user token")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}

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
	err = uh.userFollowerService.Create(ctx, userFollower)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "could not follow this user")
	}

	preloads := []string{
		"User",
		"User.Followers",
		"Follower",
		"Follower.Followers",
	}
	userFollowersToReturn, _ := uh.userFollowerService.Where(
		ctx,
		&models.UserFollower{
			BaseModel: models.BaseModel{
				ID: userFollower.ID,
			},
		},
		query.WithPreloads(preloads...),
	)
	userFollowerToReturn := &userFollowersToReturn[0]

	dtoBuilder := dto.NewDTOBuilder(uh.userService, nil)
	userDTO := dtoBuilder.BuildUserFollowerDTO(userFollowerToReturn)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, userDTO)
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
	targetID := chi.URLParam(r, "id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}

	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid auth user token")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid targed user id")
	}

	userFollowers, err := uh.userFollowerService.Where(ctx, &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	})
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "trouble with finding out if you followed to this user")
	}
	if len(userFollowers) == 0 {
		return helpers.NewAPIError(http.StatusInternalServerError, "you are not followed for this user")
	}

	userFollower := userFollowers[0]

	err = uh.userFollowerService.Delete(ctx, userFollower.ID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "can not unfollow this user")
	}

	render.Status(r, http.StatusOK)
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
	userID, _ := helpers.GetUserID(r.Context())
	userUUID, _ := uuid.Parse(userID)

	targetID := chi.URLParam(r, "id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "wrong target user id")
	}

	if userFollower, err := uh.userFollowerService.Where(context.Background(), &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	}); err != nil || len(userFollower) < 1 {
		if err != nil {
			return helpers.NewAPIError(http.StatusInternalServerError, "troubles to check if you followed")
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]bool{
			"isFollowed": false,
		})
		return nil
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]bool{
		"isFollowed": true,
	})
	return nil
}
