// chat_ws_handler.go
package chat

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/ws"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	manager        *ws.HubManager
	chatService    services.ChatService
	messageService services.MessageService
	userService    services.UserService
	dtoBuilder     *dto.DTOBuilder
	cfg            *config.Config
	logger         *slog.Logger
}

func NewChatHandler(
	manager *ws.HubManager,
	chatService services.ChatService,
	messageService services.MessageService,
	userService services.UserService,
	dtoBuilder *dto.DTOBuilder,
	cfg *config.Config,
	logger *slog.Logger,
) *ChatHandler {
	return &ChatHandler{
		manager:        manager,
		chatService:    chatService,
		messageService: messageService,
		userService:    userService,
		dtoBuilder:     dtoBuilder,
		cfg:            cfg,
		logger:         logger,
	}
}

// ServeWs handles WebSocket connections between users for private chats.
// @Summary      WebSocket connection for privat chat
// @Description  Setting WebSocket connection between authorised user and target users.
// @Tags         chat
// @Produce      json
// @Param        authToken query string true "Bearer auth token"
// @Param        userIds query string true "userIds separated by coma"
// @Param        name query string true "chan name"
// @Router       /ws/chat [get]
func (ch *ChatHandler) ServeWs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	token := r.URL.Query().Get("authToken")
	userIDRaw, err := helpers.ParseToken(ch.cfg.JwtSecret, token)
	if err != nil {
		return helpers.BadRequest("invalid auth token")
	}
	userUUID, err := uuid.Parse(userIDRaw)
	if err != nil {
		return helpers.BadRequest("invalid user ID")
	}

	userIDsParam := r.URL.Query().Get("userIds")
	if userIDsParam == "" {
		return helpers.BadRequest("userIds is required")
	}

	idStrs := strings.Split(userIDsParam, ",")
	var userIDs []uuid.UUID
	for _, idStr := range idStrs {
		id, err := uuid.Parse(strings.TrimSpace(idStr))
		if err != nil {
			return helpers.BadRequest("invalid user ID in userIds list")
		}
		userIDs = append(userIDs, id)
	}

	userIDs = append(userIDs, userUUID)
	chatName := r.URL.Query().Get("name")
	chat, err := ch.chatService.GetOrCreateGroupChat(ctx, userIDs, chatName)
	if err != nil {
		return helpers.InternalServerError("failed to get or create chat")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return helpers.InternalServerError("failed to upgrade connection")
	}

	hub, err := ch.manager.GetHub(chat.ID.String())
	if err != nil {
		return err
	}

	client := ws.NewClient(conn, hub, userUUID, chat.ID, ch.messageService, ch.userService)
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()

	// messages, err := ch.messageService.Where(ctx, &models.Message{ChatID: chat.ID}, query.WithPreloads("Sender"))
	// if err == nil {
	// 	for _, msg := range messages {
	// 		msgDTO := ch.dtoBuilder.BuildMessageDTO(&msg)
	// 		b, _ := json.Marshal(msgDTO)
	// 		client.Send <- b
	// 	}
	// }
	b, _ := json.Marshal(map[string]any{
		"type": "DH_INIT",
		"p":    hub.DhKeys.Prime,
		"q":    hub.DhKeys.Generator,
	})
	hub.Broadcast <- b

	return nil
}
