package models

type User struct {
	BaseModel

	Username        string `json:"username"`
	Email           string `gorm:"unique" json:"email"`
	PasswordHash    string `json:"passwordHash"`
	Role            string `json:"role"`
	ProfileInfo     string `json:"profileInfo"`
	ProfilePicture  string `json:"profilePictureUrl"`
	IsGoogleAccount bool   `json:"isGoogleAccount"`

	Songs           []Song           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments        []Comment        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Collections     []Collection     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Reactions       []UserReaction   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	UserCollections []UserCollection `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Chats1    []Chat         `gorm:"foreignKey:UserID1;" json:"chats1"`
	Chats2    []Chat         `gorm:"foreignKey:UserID2;" json:"chats2"`
	Messages  []Message      `gorm:"foreignKey:SenderID" json:"messages"`
	Followers []UserFollower `gorm:"foreignKey:UserID;" json:"followers"`
	Follows   []UserFollower `gorm:"foreignKey:FollowerID;" json:"follows"`
}
