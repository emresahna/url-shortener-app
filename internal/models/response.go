package models

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
	Username string     `json:"username"`
	Urls     []UserUrls `json:"urls"`
}

type UserUrls struct {
	Url         string `json:"url"`
	OriginalUrl string `json:"original_url"`
	ClickCount  int64  `json:"click_count"`
}
