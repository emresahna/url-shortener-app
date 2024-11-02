package models

type AuthHeader struct {
	Token string `header:"Authorization"`
}

type ShortenURLRequest struct {
	AuthHeader
	OriginalUrl string `json:"original_url"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
}
