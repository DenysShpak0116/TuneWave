package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
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
func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userIDuuid, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid user ID", err)
		return
	}
	user, err := uh.UserService.GetFullDTOByID(context.Background(), userIDuuid)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusNotFound, "user not found", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

type UserChatsResponse struct {
	Chats []ChatPreview `json:"chats"`
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
func (uh *UserHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(string)
	if !ok {
		fmt.Printf("token: %+v\n", r.Context().Value("userID"))
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	users, err := uh.UserService.Where(context.Background(), &models.User{
		BaseModel: models.BaseModel{
			ID: userUUID,
		},
	}, "Chats1", "Chats2", "Chats1.Messages", "Chats2.Messages",
		"Chats1.User1", "Chats1.User2", "Chats2.User1", "Chats2.User2")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusNotFound, "user not found", err)
		return
	}

	if len(users) == 0 {
		handlers.RespondWithError(w, r, http.StatusNotFound, "user not found", err)
		return
	}

	chats := make([]ChatPreview, 0)
	for _, chat := range users[0].Chats1 {
		if len(chat.Messages) == 0 {
			continue
		}

		var (
			userAvatar   string
			username     string
			targetUserID uuid.UUID
		)

		if chat.User2.ID == userUUID {
			userAvatar = chat.User1.ProfilePicture
			username = chat.User1.Username
			targetUserID = chat.User1.ID
		} else {
			userAvatar = chat.User2.ProfilePicture
			username = chat.User2.Username
			targetUserID = chat.User2.ID
		}

		chats = append(chats, ChatPreview{
			ID:           chat.ID,
			UserAvatar:   userAvatar,
			Username:     username,
			LastMessage:  chat.Messages[len(chat.Messages)-1].Content,
			TargetUserID: targetUserID,
		})
	}
	for _, chat := range users[0].Chats2 {
		if len(chat.Messages) == 0 {
			continue
		}

		var (
			userAvatar   string
			username     string
			targetUserID uuid.UUID
		)

		if chat.User1.ID == userUUID {
			userAvatar = chat.User2.ProfilePicture
			username = chat.User2.Username
			targetUserID = chat.User2.ID
		} else {
			userAvatar = chat.User1.ProfilePicture
			username = chat.User1.Username
			targetUserID = chat.User1.ID
		}

		chats = append(chats, ChatPreview{
			ID:           chat.ID,
			UserAvatar:   userAvatar,
			Username:     username,
			LastMessage:  chat.Messages[len(chat.Messages)-1].Content,
			TargetUserID: targetUserID,
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, UserChatsResponse{
		Chats: chats,
	})
}

func printChat(chat interface{}) {
	bytes, err := json.MarshalIndent(chat, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling chat: %v\n", err)
		return
	}
	fmt.Printf("Chat:\n%s\n", string(bytes))
}
