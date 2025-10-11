package models

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
)

type User struct {
	ID        string `bson:"_id"`
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
}

func NewUser(id, firstName, lastName string) User {
	return User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
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
