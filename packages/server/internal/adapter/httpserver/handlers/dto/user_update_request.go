package dto

type UserUpdateRequest struct {
	Username    string `json:"username"`
	ProfileInfo string `json:"profileInfo"`
}
