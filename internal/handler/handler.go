package handler

import (
	"github.com/EmreSahna/url-shortener-app/internal/endpoints"
	"github.com/EmreSahna/url-shortener-app/internal/middleware/bearer"
	"github.com/EmreSahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func NewHandler(s service.Service) http.Handler {
	ep := endpoints.NewEndpoints(s)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.Post("/signup", ep.SignupUserHandler)
	r.Post("/login", ep.LoginUserHandler)
	r.Get("/me/{userId}", ep.MeHandler)

	r.Group(func(r chi.Router) {
		r.Use(bearer.BearerMiddleware)
		r.Post("/shorten/url", ep.UrlShortenerHandler)
	})
	r.Get("/redirect/{code}", ep.RedirectUrlHandler)
	r.Post("/shorten/limited-url", ep.LimitedUrlShortenerHandler)

	return r
}
