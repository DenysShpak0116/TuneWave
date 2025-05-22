package models

import "github.com/google/uuid"

type UserFollower struct {
	BaseModel

	UserID uuid.UUID `json:"userId"`
	User   User      `json:"user" gorm:"constraint:OnDelete:CASCADE"`

	FollowerID uuid.UUID `json:"followerId"`
	Follower   User      `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"follower"`
}
