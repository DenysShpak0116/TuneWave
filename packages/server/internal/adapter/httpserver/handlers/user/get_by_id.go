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

	preloads := []string{
		"ChatUsers",
		"ChatUsers.Chat",
		"ChatUsers.Chat.ChatUsers",
		"ChatUsers.Chat.ChatUsers.User",
	}
	user, err := uh.userService.GetByID(ctx, userUUID, preloads...)
	if err != nil {
		return helpers.NotFound("user not found")
	}

	chats := make([]ChatPreview, 0, len(user.ChatUsers))
	chats = appendChatsForUser(ctx, chats, userUUID, user.ChatUsers, uh)

	render.JSON(w, r, chats)
	return nil
}

func appendChatsForUser(
	ctx context.Context,
	chats []ChatPreview,
	userUUID uuid.UUID,
	userChats []models.ChatUser,
	uh *UserHandler,
) []ChatPreview {
	for _, userChat := range userChats {
		chat := userChat.Chat

		lastMessage, err := uh.messageService.Last(ctx, &models.Message{ChatID: chat.ID})
		if err != nil {
			lastMessage = &models.Message{Content: ""}
		}

		var otherUsers []models.User
		for _, cu := range chat.ChatUsers {
			if cu.UserID != userUUID {
				otherUsers = append(otherUsers, cu.User)
			}
		}

		chatName := chat.Name
		var avatar string
		if len(otherUsers) > 0 {
			if chatName == "" {
				chatName = otherUsers[0].Username
			}
			avatar = otherUsers[0].ProfilePicture
		}

		chatPreview := ChatPreview{
			ID:          chat.ID,
			UserAvatar:  avatar,
			Username:    chatName,
			LastMessage: lastMessage.Content,
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
	}
	user, err := uh.userService.GetByID(ctx, userUUID, preloads...)
	if err != nil {
		return helpers.InternalServerError("could not find user")
	}

	collections := make([]dto.CollectionDTO, 0, len(user.UserCollections))
	for _, userCollection := range user.UserCollections {
		collections = append(collections, *uh.dtoBuilder.BuildCollectionDTO(&userCollection.Collection))
	}
	render.JSON(w, r, collections)
	return nil
}
