package models

import "net/http"

type Error struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}

func UrlNotFoundErr() *Error {
	return &Error{
		Message:    "Url not found.",
		StatusCode: http.StatusNotFound,
	}
}
