package models

import "github.com/google/uuid"

type UserCollection struct {
	BaseModel

	UserID uuid.UUID `json:"userId"`
	User   User      `json:"user"`

	CollectionID uuid.UUID  `json:"collectionId"`
	Collection   Collection `json:"collection"`
}
