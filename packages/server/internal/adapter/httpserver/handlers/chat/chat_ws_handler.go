// chat_ws_handler.go
package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/ws"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/golang-jwt/jwt/v5"
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
}

func NewChatHandler(
	manager *ws.HubManager,
	chatService services.ChatService,
	messageService services.MessageService,
) *ChatHandler {
	return &ChatHandler{
		Manager:        manager,
		ChatService:    chatService,
		MessageService: messageService,
	}
}

func (ch *ChatHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVE WS WSWS\n")
	targetID := r.URL.Query().Get("targetUserId")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid target user ID", err)
		return
	}

	claimsRaw := r.Context().Value("claims")
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Invalid claims", nil)
		return
	}

	userIDRaw, ok := claims["userId"].(string)
	if !ok {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userIDRaw)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	chat, err := ch.ChatService.GetOrCreatePrivateChat(r.Context(), userUUID, targetUUID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get or create chat", err)
		return
	}

	log.Printf("Chat ID: %s", chat.ID.String())

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to upgrade connection", err)
		return
	}

	log.Printf("Connection upgraded to WebSocket")

	chatIDStr := chat.ID.String()
	hub := ch.Manager.GetHub(chatIDStr)

	client := ws.NewClient(conn, hub, userUUID, chat.ID, ch.MessageService)

	log.Printf("Client created: %s", client.UserID.String())

	hub.Register <- client
	go client.WritePump()
	go client.ReadPump()

	messages, err := ch.MessageService.Where(r.Context(), &models.Message{ChatID: chat.ID})
	if err == nil {
		for _, msg := range messages {
			b, _ := json.Marshal(msg)
			client.Send <- b
		}
	}
}
