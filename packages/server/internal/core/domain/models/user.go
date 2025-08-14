package models

type User struct {
	BaseModel

	ProfilePicture  string `json:"profilePictureUrl"`
	PasswordHash    string `json:"passwordHash"`
	ProfileInfo     string `json:"profileInfo"`
	Username        string `json:"username"`
	Email           string `gorm:"unique" json:"email"`
	Role            string `json:"role"`
	IsGoogleAccount bool   `json:"isGoogleAccount"`

	Follows         []UserFollower   `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"follows"`
   Chats2          []Chat           `gorm:"foreignKey:UserID2;constraint:OnDelete:CASCADE" json:"chats2"`
	Chats1          []Chat           `gorm:"foreignKey:UserID1;constraint:OnDelete:CASCADE" json:"chats1"`
	Followers       []UserFollower   `gorm:"constraint:OnDelete:CASCADE" json:"followers"`
	UserCollections []UserCollection `gorm:"constraint:OnDelete:CASCADE"`
	Reactions       []UserReaction   `gorm:"constraint:OnDelete:CASCADE"`
	Collections     []Collection     `gorm:"constraint:OnDelete:CASCADE"`
	Comments        []Comment        `gorm:"constraint:OnDelete:CASCADE"`
	Songs           []Song           `gorm:"constraint:OnDelete:CASCADE"`
}
