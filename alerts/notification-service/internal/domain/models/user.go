package models

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
)

type User struct {
	ID        string `bson:"_id"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Role      string `bson:"role"`
}

func NewPartialUser(id string) *User {
	return &User{
		ID:        id,
		FirstName: "",
		LastName:  "",
	}
}

func NewUser(id, firstName, lastName, role string) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	}
}

func (u *User) DTO() dto.UserDTO {
	return dto.UserDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
