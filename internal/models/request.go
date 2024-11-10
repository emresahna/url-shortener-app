package models

type SignupUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHeader struct {
	Token string `json:"-"`
}

type ShortenURLRequest struct {
	AuthHeader
	OriginalUrl string `json:"original_url"`
}
