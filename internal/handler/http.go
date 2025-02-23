package handler

import (
	"net/http"

	"github.com/emresahna/url-shortener-app/configs"
	"github.com/emresahna/url-shortener-app/internal/endpoints"
	"github.com/emresahna/url-shortener-app/internal/middleware/bearer"
	"github.com/emresahna/url-shortener-app/internal/middleware/ipaddr"
	"github.com/emresahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewHTTP(s service.Service, cfg configs.Cors) http.Handler {
	ep := endpoints.New(s)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: cfg.AllowedMethods,
		AllowedHeaders: cfg.AllowedHeaders,
	}))

	r.Post("/user/signup", ep.UserSignupHandler)
	r.Post("/user/login", ep.UserLoginHandler)
	r.Post("/token/refresh", ep.TokenRefreshHandler)
	r.Get("/url/redirect/{code}", ep.UrlRedirectHandler)

	r.Group(func(r chi.Router) {
		r.Use(ipaddr.Middleware)
		r.Use(bearer.Middleware)
		r.Get("/user/me", ep.UserMeHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(bearer.Middleware)
		r.Post("/url/shorten/user", ep.UrlShortenUserHandler)
		r.Delete("/url/remove/{id}", ep.UrlRemoveHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(ipaddr.Middleware)
		r.Post("/url/shorten/guest", ep.UrlShortenGuestHandler)
	})

	return r
}
