package models

type Chat struct {
	BaseModel
	Name      string     `json:"name"`
	ChatUsers []ChatUser `gorm:"constraint:OnDelete:CASCADE"`
}
