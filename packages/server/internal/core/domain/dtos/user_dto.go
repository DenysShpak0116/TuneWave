package dtos

import "github.com/google/uuid"

type UserDTO struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	ProfileInfo string    `json:"profileInfo"`
	Email       string    `json:"email"`
	ProfilePic  string    `json:"profilePictureUrl"`
}
