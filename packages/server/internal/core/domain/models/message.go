package models

import "github.com/google/uuid"

type Message struct {
	BaseModel

	Content string `json:"content"`

	ChatID   uuid.UUID `json:"chatId"`
	Chat     Chat      `gorm:"constraint:OnDelete:CASCADE"`
	SenderID uuid.UUID `json:"senderId"`
	Sender   User      `gorm:"constraint:OnDelete:CASCADE"`
}
