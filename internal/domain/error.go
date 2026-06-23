package domain

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var ErrUserAlreadyExists = errors.New("User already exists")
var ErrUserNotFound = errors.New("User not found")

type Error struct {
	Time time.Time
	Message string
}

func FormError(message string, time time.Time, w http.ResponseWriter, code int) error {
	error := Error{
		Message: message,
		Time: time,
	}
	b, err := json.MarshalIndent(error, "", "	")
	if err != nil { return err }
	http.Error(w, string(b), code)
	return nil
}