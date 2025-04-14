package models

import "github.com/google/uuid"

type Message struct {
	BaseModel

	Content string `json:"content"`

	ChatID uuid.UUID `gorm:"type:uuid" json:"chatId"`
	Chat   Chat      `gorm:"foreignKey:ChatID"`

	SenderID uuid.UUID `gorm:"type:uuid" json:"senderId"`
	Sender   User      `gorm:"foreignKey:SenderID"`
}
