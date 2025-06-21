// chat_ws_handler.go
package chat

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/ws"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	Manager        *ws.HubManager
	ChatService    services.ChatService
	MessageService services.MessageService
	JWTSecret      string
}

func NewChatHandler(
	manager *ws.HubManager,
	chatService services.ChatService,
	messageService services.MessageService,
	cfg *config.Config,
) *ChatHandler {
	return &ChatHandler{
		Manager:        manager,
		ChatService:    chatService,
		MessageService: messageService,
		JWTSecret:      cfg.JwtSecret,
	}
}

// ServeWs handles WebSocket connections between users for private chats.
// @Summary      WebSocket connection for privat chat
// @Description  Setting WebSocket connection between authorised user and target user by `targetUserId`.
// @Tags         chat
// @Produce      json
// @Param        targetUserId query string true "UUID of target user"
// @Param        authToken query string true "Bearer auth token"
// @Router       /ws/chat [get]
func (ch *ChatHandler) ServeWs(w http.ResponseWriter, r *http.Request) error {
	targetID := r.URL.Query().Get("targetUserId")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid target user ID")
	}

	token := r.URL.Query().Get("authToken")
	userIDRaw, err := helpers.ParseToken(ch.JWTSecret, token)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid auth token")
	}

	userUUID, err := uuid.Parse(userIDRaw)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}

	chat, err := ch.ChatService.GetOrCreatePrivateChat(r.Context(), userUUID, targetUUID)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get or create chat")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to upgrade connection")
	}

	chatIDStr := chat.ID.String()
	hub := ch.Manager.GetHub(chatIDStr)

	client := ws.NewClient(conn, hub, userUUID, chat.ID, ch.MessageService)

	hub.Register <- client
	go client.WritePump()
	go client.ReadPump()

	messages, err := ch.MessageService.Where(r.Context(), &models.Message{ChatID: chat.ID})
	dtoBuilder := dto.NewDTOBuilder(nil, nil)
	if err == nil {
		for _, msg := range messages {
			msgDTO := dtoBuilder.BuildMessageDTO(&msg)
			b, _ := json.Marshal(msgDTO)
			client.Send <- b
		}
	}

	return nil
}
