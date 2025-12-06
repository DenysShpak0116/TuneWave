package models

import "github.com/google/uuid"

type ChatUser struct {
	BaseModel

	ChatID uuid.UUID `json:"chatId"`
	Chat   Chat      `gorm:"foreignKey:ChatID" json:"chat"`

	UserID uuid.UUID `json:"userId"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
}
