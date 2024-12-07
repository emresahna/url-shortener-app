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

func UrlNotValidErr() *Error {
	return &Error{
		Message:    "Url not valid.",
		StatusCode: http.StatusBadRequest,
	}
}

func InternalServerErr() *Error {
	return &Error{
		Message:    "Something went wrong.",
		StatusCode: http.StatusInternalServerError,
	}
}

func TokenFailureErr() *Error {
	return &Error{
		Message:    "Credentials could not validate.",
		StatusCode: http.StatusUnauthorized,
	}
}

func TokenExpiredErr() *Error {
	return &Error{
		Message:    "Token expired.",
		StatusCode: http.StatusUnauthorized,
	}
}

func CustomerIdParseErr() *Error {
	return &Error{
		Message:    "Customer ID couldn't validate.",
		StatusCode: http.StatusUnauthorized,
	}
}

func CreateURLErr() *Error {
	return &Error{
		Message:    "Error while saving url to database.",
		StatusCode: http.StatusBadRequest,
	}
}

func CreateClickCountErr() *Error {
	return &Error{
		Message:    "Error while creating click count record to database.",
		StatusCode: http.StatusBadRequest,
	}
}

func SaveToCacheErr() *Error {
	return &Error{
		Message:    "Error while caching url.",
		StatusCode: http.StatusBadRequest,
	}
}

func ExpireTimeAlreadyPassedErr() *Error {
	return &Error{
		Message:    "Expire time already passed.",
		StatusCode: http.StatusBadRequest,
	}
}

func ParseExpireTimeErr() *Error {
	return &Error{
		Message:    "Error while parsing expire time.",
		StatusCode: http.StatusBadRequest,
	}
}

func UserNotFoundErr() *Error {
	return &Error{
		Message:    "User not found.",
		StatusCode: http.StatusNotFound,
	}
}

func PasswordDoesntMatchErr() *Error {
	return &Error{
		Message:    "Password does not match with account.",
		StatusCode: http.StatusBadRequest,
	}
}

func ErrorWhileCreatingTokenErr() *Error {
	return &Error{
		Message:    "Error while creating token.",
		StatusCode: http.StatusBadRequest,
	}
}

func ErrorWhileEncryptingPasswordErr() *Error {
	return &Error{
		Message:    "Error while encrypting password.",
		StatusCode: http.StatusBadRequest,
	}
}

func ErrorWhileCreatingUserErr() *Error {
	return &Error{
		Message:    "Error while creating user.",
		StatusCode: http.StatusBadRequest,
	}
}

func UserAlreadyExistsErr() *Error {
	return &Error{
		Message:    "User already exists.",
		StatusCode: http.StatusBadRequest,
	}
}

func FreeTierExceedErr() *Error {
	return &Error{
		Message:    "Your free tier usage ended. Please try later...",
		StatusCode: http.StatusBadRequest,
	}
}
