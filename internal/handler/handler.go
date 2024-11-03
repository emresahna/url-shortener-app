package handler

import (
	"github.com/EmreSahna/url-shortener-app/internal/endpoints"
	"github.com/EmreSahna/url-shortener-app/internal/middleware"
	"github.com/EmreSahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewHandler(s service.Service) http.Handler {
	ep := endpoints.NewEndpoints(s)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.BearerMiddleware)
		r.Post("/url-shortener", ep.UrlShortenerHandler)
	})

	r.Get("/get-user/{userID}", ep.GetUserHandler)
	r.Post("/create-user", ep.CreateUserHandler)

	return r
}
