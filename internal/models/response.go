package models

import "github.com/google/uuid"

type SignupUserResponse struct {
	Info string `json:"info"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ShortenURLResponse struct {
	Url string `json:"url"`
}

type UserResponse struct {
	Username string     `json:"username,omitempty"`
	Urls     []UserUrls `json:"urls,omitempty"`
}

type UserUrls struct {
	Id          uuid.UUID `json:"id"`
	Url         string    `json:"url"`
	OriginalUrl string    `json:"original_url"`
	ClickCount  int64     `json:"click_count"`
	IsActive    *bool     `json:"is_active"`
	IsDeleted   *bool     `json:"is_deleted"`
}

type RemoveUrlResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}
