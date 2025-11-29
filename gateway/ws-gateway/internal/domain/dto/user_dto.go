package dto

import "fmt"

type UserDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *UserDTO) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
