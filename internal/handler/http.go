// Package handler provides HTTP request handling and routing
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/emresahna/url-shortener-app/configs"
	_ "github.com/emresahna/url-shortener-app/internal/docs"
	"github.com/emresahna/url-shortener-app/internal/endpoints"
	"github.com/emresahna/url-shortener-app/internal/middleware/bearer"
	"github.com/emresahna/url-shortener-app/internal/middleware/ipaddr"
	"github.com/emresahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHTTP(s service.Service, cfg configs.Cors) http.Handler {
	ep := endpoints.New(s)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: cfg.AllowedMethods,
		AllowedHeaders: cfg.AllowedHeaders,
	}))

	// Add request timestamp to context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "requestTime", time.Now().Format(time.RFC3339))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// Health check endpoints
	r.Get("/health", ep.HealthCheckHandler)
	r.Get("/health/ready", ep.ReadinessCheckHandler)
	r.Get("/health/live", ep.LivenessCheckHandler)

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
	))

	// API versioning with v1 prefix
	r.Route("/api/v1", func(r chi.Router) {
		// Public endpoints
		r.Post("/user/signup", ep.UserSignupHandler)
		r.Post("/user/login", ep.UserLoginHandler)
		r.Post("/token/refresh", ep.TokenRefreshHandler)
		r.Get("/url/redirect/{code}", ep.UrlRedirectHandler)

		// User endpoints with IP and bearer token middleware
		r.Group(func(r chi.Router) {
			r.Use(ipaddr.Middleware)
			r.Use(bearer.Middleware)
			r.Get("/user/me", ep.UserMeHandler)
		})

		// URL endpoints with bearer token middleware
		r.Group(func(r chi.Router) {
			r.Use(bearer.Middleware)
			r.Post("/url/shorten/user", ep.UrlShortenUserHandler)
			r.Delete("/url/remove/{id}", ep.UrlRemoveHandler)
		})

		// Guest URL endpoints with IP middleware
		r.Group(func(r chi.Router) {
			r.Use(ipaddr.Middleware)
			r.Post("/url/shorten/guest", ep.UrlShortenGuestHandler)
		})
	})

	return r
}
