package domain

import (
	"errors"
	"strings"
)

type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func (u *User) ValidateUser() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("Username or password is empty")
	}
	if len(u.Username) > 20 {
		return errors.New("Username is too long, max length is 20")
	}
	if len(u.Password) > 10 {
		return errors.New("Password is too long, max length is 10")
	}
	if strings.Contains(u.Username, " ") {
		return errors.New("Username must be entered without spaces")
	}
	if strings.Contains(u.Password, " ") {
		return errors.New("Password must be entered without spaces")
	}
	return nil
}