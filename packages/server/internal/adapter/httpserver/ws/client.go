package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn           *websocket.Conn
	Send           chan []byte
	Hub            *Hub
	MessageService services.MessageService
	UserID         uuid.UUID
	ChatID         uuid.UUID
}

func NewClient(conn *websocket.Conn, hub *Hub, userID, chatID uuid.UUID, messageService services.MessageService) *Client {
	return &Client{
		Conn:           conn,
		Send:           make(chan []byte),
		Hub:            hub,
		MessageService: messageService,
		UserID:         userID,
		ChatID:         chatID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ReadPump] panic: %v", r)
		}
		c.Hub.Unregister <- c
		c.Conn.Close()
		log.Println("[ReadPump] client disconnected")
	}()

	dtoBuilder := dto.NewDTOBuilder()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("[ReadPump] read error:", err)
			break
		}
		log.Printf("[ReadPump] Received: %s", msg)

		var payload struct {
			Content string `json:"content"`
		}

		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("[ReadPump] invalid format:", err)
			continue
		}
		log.Printf("[ReadPump] Parsed content: %s", payload.Content)

		message := &models.Message{
			Content:  payload.Content,
			ChatID:   c.ChatID,
			SenderID: c.UserID,
		}

		if err := c.MessageService.Create(context.Background(), message); err != nil {
			continue
		}

		messageDTO := dtoBuilder.BuildMessageDTO(message)
		outgoing, _ := json.Marshal(messageDTO)
		c.Hub.Broadcast <- outgoing
	}
}

func (c *Client) WritePump() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WritePump] panic: %v", r)
		}
		c.Conn.Close()
	}()

	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
