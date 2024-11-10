package models

import "time"

type SignupUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ShortenURLRequest struct {
	OriginalUrl string     `json:"original_url"`
	ExpireTime  *time.Time `json:"expire_time"`
}
