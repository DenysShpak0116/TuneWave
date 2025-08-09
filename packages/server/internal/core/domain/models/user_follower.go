package models

import "github.com/google/uuid"

type UserFollower struct {
	BaseModel

	UserID uuid.UUID `json:"userId"`
	User   User      ` gorm:"foreignKey:UserID" json:"user"`

	FollowerID uuid.UUID `json:"followerId"`
	Follower   User      `gorm:"foreignKey:FollowerID" json:"follower"`
}
