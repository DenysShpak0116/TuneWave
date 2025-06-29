package user

import (
	"context"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GetByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Router /users/{id} [get]
func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("invalid user ID")
	}

	preloads := []string{"Follows", "Follows.User", "Followers", "Followers.Follower"}
	user, err := uh.userService.GetByID(ctx, userUUID, preloads...)
	if err != nil {
		return helpers.NotFound("user not found")
	}

	render.JSON(w, r, uh.dtoBuilder.BuildFullUserDTO(user))
	return nil
}

type ChatPreview struct {
	ID           uuid.UUID `json:"id"`
	TargetUserID uuid.UUID `json:"targetUserId"`
	UserAvatar   string    `json:"userAvatar"`
	Username     string    `json:"username"`
	LastMessage  string    `json:"lastMessage"`
}

// GetChats godoc
// @Summary Get user chats
// @Description Get user chats
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Router /chats [get]
func (uh *UserHandler) GetChats(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userUUID, err := helpers.GetUserID(ctx)
	if err != nil {
		return helpers.NotFound("invalid user ID")
	}

	preloads := []string{"Chats1", "Chats2", "Chats1.User1", "Chats1.User2", "Chats2.User1", "Chats2.User2"}
	user, err := uh.userService.GetByID(ctx, userUUID, preloads...)
	if err != nil {
		return helpers.NotFound("user not found")
	}

	chats := make([]ChatPreview, 0)
	chats = appendChatsForUser(ctx, chats, userUUID, user.Chats1, uh)
	chats = appendChatsForUser(ctx, chats, userUUID, user.Chats2, uh)
	render.JSON(w, r, chats)
	return nil
}

func appendChatsForUser(ctx context.Context, chats []ChatPreview, userUUID uuid.UUID, userChats []models.Chat, uh *UserHandler) []ChatPreview {
	for _, chat := range userChats {
		lastMessage, err := uh.messageService.Last(ctx, &models.Message{ChatID: chat.ID})
		if err != nil {
			lastMessage = &models.Message{
				Content: "",
			}
		}

		var chatPreview ChatPreview
		if chat.User1.ID == userUUID {
			chatPreview = ChatPreview{
				ID:           chat.ID,
				UserAvatar:   chat.User2.ProfilePicture,
				Username:     chat.User2.Username,
				LastMessage:  lastMessage.Content,
				TargetUserID: chat.User2.ID,
			}
		} else {
			chatPreview = ChatPreview{
				ID:           chat.ID,
				UserAvatar:   chat.User1.ProfilePicture,
				Username:     chat.User1.Username,
				LastMessage:  lastMessage.Content,
				TargetUserID: chat.User1.ID,
			}
		}

		chats = append(chats, chatPreview)
	}
	return chats
}

// GetUserCollections godoc
// @Tags user
// @Produce json
// @Param id path string true "User id"
// @Router /users/{id}/collections [get]
func (uh *UserHandler) GetUserCollections(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return helpers.BadRequest("wrong user id")
	}

	preloads := []string{
		"UserCollections",
		"UserCollections.Collection",
		"UserCollections.Collection.User",
		"UserCollections.Collection.User.Followers",
	}
	user, err := uh.userService.GetByID(ctx, userUUID, preloads...)
	if err != nil {
		return helpers.InternalServerError("could not find user")
	}

	collections := make([]dto.CollectionDTO, 0)
	for _, userCollection := range user.UserCollections {
		collections = append(collections, *uh.dtoBuilder.BuildCollectionDTO(&userCollection.Collection))
	}
	render.JSON(w, r, collections)
	return nil
}
