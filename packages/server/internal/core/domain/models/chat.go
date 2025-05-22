package models

import "github.com/google/uuid"

type Chat struct {
	BaseModel

	UserID1 uuid.UUID `gorm:"not null" json:"userId1"`
	User1   User      `gorm:"foreignKey:UserID1;" json:"user1"`

	UserID2 uuid.UUID `gorm:"not null" json:"userId2"`
	User2   User      `gorm:"foreignKey:UserID2;" json:"user2"`

	Messages []Message `json:"messages"`
}
