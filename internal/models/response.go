package models

type SignupUserResponse struct {
	Info string `json:"inasd"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ShortenURLResponse struct {
	Url string `json:"url"`
}
