package service

import (
	"context"

	"github.com/emresahna/url-shortener-app/internal/auth"
	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/redis"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RedirectionsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "url_shortener_redirections_total",
		Help: "The total number of URL redirections",
	})
	ShorteningsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "url_shortener_shortenings_total",
		Help: "The total number of URL shortenings",
	}, []string{"type"})
)

// Service defines the complete service interface
type Service interface {
	HealthCheck(ctx context.Context) (models.HealthResponse, error)
	ReadinessCheck(ctx context.Context) (models.HealthResponse, error)
	LivenessCheck(ctx context.Context) (models.HealthResponse, error)
	UrlRedirect(ctx context.Context, code string) (string, error)
	UrlRemove(ctx context.Context, id string) (models.RemoveUrlResponse, error)
	UrlShortenGuest(
		ctx context.Context,
		req models.ShortenURLRequest,
	) (models.ShortenURLResponse, error)
	UrlShortenUser(
		ctx context.Context,
		req models.ShortenURLRequest,
	) (models.ShortenURLResponse, error)
	UserLogin(ctx context.Context, req models.LoginUserRequest) (models.LoginUserResponse, error)
	UserSignup(ctx context.Context, req models.SignupUserRequest) (models.SignupUserResponse, error)
	UserMe(ctx context.Context) (models.UserResponse, error)
	TokenRefresh(
		ctx context.Context,
		req models.RefreshTokenRequest,
	) (models.LoginUserResponse, error)
}

// service implements the auth Service interface
type service struct {
	db  sqlc.Querier
	rcc redis.Store
	jwt auth.Auth
	rca redis.Store
}

// NewService creates a new auth service
func New(db sqlc.Querier, rcc redis.Store, jwt auth.Auth, rca redis.Store) Service {
	return &service{
		db:  db,
		rcc: rcc,
		jwt: jwt,
		rca: rca,
	}
}
