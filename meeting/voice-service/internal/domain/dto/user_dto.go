package dto

type UserDTO struct {
	ID        string  `json:"id"`
	AvatarURL *string `json:"avatarURL,omitempty"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
}
