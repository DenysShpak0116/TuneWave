package models

type User struct {
	BaseModel

	Username        string `json:"username"`
	Email           string `gorm:"unique" json:"email"`
	PasswordHash    string `json:"passwordHash"`
	ProfileInfo     string `json:"profileInfo"`
	ProfilePicture  string `json:"profilePictureUrl"`
	IsGoogleAccount bool   `json:"isGoogleAccount"`
	Songs           []Song
	Collections     []Collection
	Comments        []Comment
	Reactions       []UserReaction

	Chats1 []Chat `gorm:"foreignKey:UserID1;references:ID" json:"chats1"`
	Chats2 []Chat `gorm:"foreignKey:UserID2;references:ID" json:"chats2"`

	Messages []Message `gorm:"foreignKey:SenderID" json:"messages"`
}
