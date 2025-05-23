package user

import (
	"context"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
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
func (us *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid targed user id", err)
		return
	}

	userID, err := helpers.GetUserID(r.Context())
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid auth user token", err)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid targed user id", err)
		return
	}

	userFollowers, _ := us.UserFollowerService.Where(context.Background(), &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	})
	if len(userFollowers) > 0 {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "you are already followed to this user", err)
		return
	}

	userFollower := &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	}
	err = us.UserFollowerService.Create(context.Background(), userFollower)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "could not follow this user", err)
		return
	}

	userFollowersToReturn, _ := us.UserFollowerService.Where(context.Background(), &models.UserFollower{
		BaseModel: models.BaseModel{
			ID: userFollower.ID,
		},
	}, "User", "User.Followers", "Follower", "Follower.Followers")
	userFollowerToReturn := userFollowersToReturn[0]

	userDTO := dtos.UserFollowerDTO{
		ID: userFollowerToReturn.ID,
		User: dtos.UserDTO{
			ID:             userFollowerToReturn.User.ID,
			Username:       userFollowerToReturn.User.Username,
			Role:           userFollowerToReturn.User.Role,
			ProfilePicture: userFollowerToReturn.User.ProfilePicture,
			ProfileInfo:    userFollowerToReturn.User.ProfileInfo,
			Followers:      int64(len(userFollowerToReturn.User.Followers)),
		},
		Follower: dtos.UserDTO{
			ID:             userFollowerToReturn.Follower.ID,
			Username:       userFollowerToReturn.Follower.Username,
			Role:           userFollowerToReturn.Follower.Role,
			ProfilePicture: userFollowerToReturn.Follower.ProfilePicture,
			ProfileInfo:    userFollowerToReturn.Follower.ProfileInfo,
			Followers:      int64(len(userFollowerToReturn.Follower.Followers)),
		},
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, userDTO)
}

// UnfollowUser	godoc
// @Summary		Unfollow user using id param
// @Description	Unfollow user using id param
// @Tags		user
// @Security	BearerAuth
// @Produce		json
// @Param		id path string true "User to unfollow ID"
// @Router		/users/{id}/unfollow [delete]
func (uh *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid targed user id", err)
		return
	}

	userID, err := helpers.GetUserID(r.Context())
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid auth user token", err)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid targed user id", err)
		return
	}

	userFollowers, err := uh.UserFollowerService.Where(context.Background(), &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "trouble with finding out if you followed to this user", err)
		return
	}
	if len(userFollowers) == 0 {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "you are not followed for this user", err)
		return
	}

	userFollower := userFollowers[0]

	err = uh.UserFollowerService.Delete(context.Background(), userFollower.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "can not unfollow this user", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.NoContent(w, r)
}

// IsFollowed godoc
// @Tags		user
// @Security	BearerAuth
// @Produce json
// @Param id path string true "The user ID you use to check if you are followed"
// @Router /users/{id}/is-followed [get]
func (uh *UserHandler) IsFollowed(w http.ResponseWriter, r *http.Request) {
	userID, _ := helpers.GetUserID(r.Context())
	userUUID, _ := uuid.Parse(userID)

	targetID := chi.URLParam(r, "id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "wrong target user id", err)
		return
	}

	if userFollower, err := uh.UserFollowerService.Where(context.Background(), &models.UserFollower{
		UserID:     targetUUID,
		FollowerID: userUUID,
	}); err != nil || len(userFollower) < 1 {
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "troubles to check if you followed", err)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]bool{
			"isFollowed": false,
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]bool{
		"isFollowed": true,
	})
}
