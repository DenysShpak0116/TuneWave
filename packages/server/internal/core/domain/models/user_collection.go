package models

import "github.com/google/uuid"

type UserCollection struct {
	BaseModel

	UserID uuid.UUID `json:"userId"`
	User   User      `json:"user" gorm:"ForeignKey:UserID;constraint:OnDelete:CASCADE"`

	CollectionID uuid.UUID  `json:"collectionId"`
	Collection   Collection `json:"collection" gorm:"ForeignKey:CollectionID;constraint:OnDelete:CASCADE"`
}
