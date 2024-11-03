package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ShortenURLResponse struct {
	Url string `json:"url"`
}

type CreateUserResponse struct {
	Username string `json:"username"`
}

type GetUserResponse struct {
	Id        uuid.UUID          `json:"id"`
	Username  string             `json:"username"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
