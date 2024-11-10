package handler

import (
	"github.com/EmreSahna/url-shortener-app/internal/endpoints"
	"github.com/EmreSahna/url-shortener-app/internal/middleware/bearer"
	"github.com/EmreSahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewHandler(s service.Service) http.Handler {
	ep := endpoints.NewEndpoints(s)

	r := chi.NewRouter()

	r.Post("/signup", ep.SignupUserHandler)
	r.Post("/login", ep.LoginUserHandler)
	r.Group(func(r chi.Router) {
		r.Use(bearer.BearerMiddleware)
		r.Post("/url-shortener", ep.UrlShortenerHandler)
	})

	return r
}
