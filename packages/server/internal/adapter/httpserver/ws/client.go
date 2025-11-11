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
	UserService    services.UserService
	UserID         uuid.UUID
	ChatID         uuid.UUID
}

func NewClient(
	conn *websocket.Conn,
	hub *Hub,
	userID, chatID uuid.UUID,
	messageService services.MessageService,
	userService services.UserService,
) *Client {
	return &Client{
		Conn:           conn,
		Send:           make(chan []byte),
		Hub:            hub,
		MessageService: messageService,
		UserService:    userService,
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

	dtoBuilder := dto.NewDTOBuilder(c.UserService, nil)

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("[ReadPump] read error:", err)
			break
		}

		var payload struct {
			Receiver string `json:"receiver"`
			Content  string `json:"content"`
		}
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("[ReadPump] invalid format:", err)
			continue
		}

		var recieverID uuid.UUID
		if payload.Receiver != "" {
			recieverID, _ = uuid.Parse(payload.Receiver)
		}

		user, err := c.UserService.GetByID(context.TODO(), c.UserID)
		if err != nil {
			log.Println("Failed to get user by id:", c.UserID)
			return
		}

		message := &models.Message{
			Content:  payload.Content,
			ChatID:   c.ChatID,
			SenderID: c.UserID,
		}

		if err := c.MessageService.Create(context.Background(), message); err != nil {
			continue
		}

		message.Sender = *user

		messageDTO := dtoBuilder.BuildMessageDTO(message)
		outgoing, _ := json.Marshal(messageDTO)
		if len(recieverID) == 0 {
			c.Hub.Broadcast <- outgoing
		} else {
			for cl := range c.Hub.Clients {
				if cl.UserID == recieverID {
					if !c.Hub.Clients[cl] {
						break
					}
					cl.Send <- []byte(payload.Content)
				}
			}
		}
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
