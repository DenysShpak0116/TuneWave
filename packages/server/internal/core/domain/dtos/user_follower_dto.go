package dtos

import "github.com/google/uuid"

type UserFollowerDTO struct {
	ID       uuid.UUID `json:"id"`
	User     UserDTO   `json:"user"`
	Follower UserDTO   `json:"follower"`
}
