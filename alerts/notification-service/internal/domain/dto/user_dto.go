package dto

import "fmt"

type UserDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (u *UserDTO) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
