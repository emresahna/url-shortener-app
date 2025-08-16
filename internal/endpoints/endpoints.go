package endpoints

import (
	"net/http"

	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
)

type Endpoints interface {
	TokenRefreshHandler(w http.ResponseWriter, r *http.Request)
	UserMeHandler(w http.ResponseWriter, r *http.Request)
	UserSignupHandler(w http.ResponseWriter, r *http.Request)
	UserLoginHandler(w http.ResponseWriter, r *http.Request)
	UrlShortenUserHandler(w http.ResponseWriter, r *http.Request)
	UrlRemoveHandler(w http.ResponseWriter, r *http.Request)
	UrlShortenGuestHandler(w http.ResponseWriter, r *http.Request)
	UrlRedirectHandler(w http.ResponseWriter, r *http.Request)
	HealthCheckHandler(w http.ResponseWriter, r *http.Request)
	ReadinessCheckHandler(w http.ResponseWriter, r *http.Request)
	LivenessCheckHandler(w http.ResponseWriter, r *http.Request)
}

type endpoints struct {
	s service.Service
}

func New(s service.Service) Endpoints {
	return &endpoints{
		s: s,
	}
}

// HealthCheckHandler godoc
// @Summary Health check endpoint
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /health [get]
func (e *endpoints) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := e.s.HealthCheck(r.Context())
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, http.StatusOK)
}

// ReadinessCheckHandler godoc
// @Summary Readiness check endpoint
// @Description Check if the service is ready to accept requests
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /health/ready [get]
func (e *endpoints) ReadinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := e.s.ReadinessCheck(r.Context())
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, http.StatusOK)
}

// LivenessCheckHandler godoc
// @Summary Liveness check endpoint
// @Description Check if the service is alive
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /health/live [get]
func (e *endpoints) LivenessCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := e.s.LivenessCheck(r.Context())
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, http.StatusOK)
}

// TokenRefreshHandler godoc
// @Summary Refresh access token
// @Description Refresh an expired access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.LoginUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (e *endpoints) TokenRefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	resp, err := e.s.TokenRefresh(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}

// UserMeHandler godoc
// @Summary Get current user information
// @Description Get the current authenticated user's profile and URLs
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/user/me [get]
func (e *endpoints) UserMeHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := e.s.UserMe(r.Context())
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}

// UserSignupHandler godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.SignupUserRequest true "User signup request"
// @Success 201 {object} models.SignupUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /api/v1/user/signup [post]
func (e *endpoints) UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SignupUserRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	resp, err := e.s.UserSignup(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 201)
}

// UserLoginHandler godoc
// @Summary User login
// @Description Authenticate user and return access/refresh tokens
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.LoginUserRequest true "User login request"
// @Success 200 {object} models.LoginUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/user/login [post]
func (e *endpoints) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	resp, err := e.s.UserLogin(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}

// UrlShortenUserHandler godoc
// @Summary Shorten URL for authenticated user
// @Description Create a shortened URL for an authenticated user
// @Tags url
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ShortenURLRequest true "URL shorten request"
// @Success 201 {object} models.ShortenURLResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/url/shorten [post]
func (e *endpoints) UrlShortenUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenURLRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	resp, err := e.s.UrlShortenUser(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 201)
}

// UrlRemoveHandler godoc
// @Summary Remove/delete a shortened URL
// @Description Delete a shortened URL by ID (only owner can delete)
// @Tags url
// @Produce json
// @Security BearerAuth
// @Param id path string true "URL ID"
// @Success 200 {object} models.RemoveUrlResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/url/{id} [delete]
func (e *endpoints) UrlRemoveHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := e.s.UrlRemove(r.Context(), id)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}

// UrlShortenGuestHandler godoc
// @Summary Shorten URL for guest user
// @Description Create a shortened URL without authentication (guest mode)
// @Tags url
// @Accept json
// @Produce json
// @Param request body models.ShortenURLRequest true "URL shorten request"
// @Success 201 {object} models.ShortenURLResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/url/shorten/guest [post]
func (e *endpoints) UrlShortenGuestHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenURLRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	resp, err := e.s.UrlShortenGuest(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 201)
}

// UrlRedirectHandler godoc
// @Summary Redirect to original URL
// @Description Redirect to the original URL using the shortened code
// @Tags url
// @Param code path string true "Shortened URL code"
// @Success 302 "Redirect to original URL"
// @Failure 404 {object} models.ErrorResponse
// @Failure 410 {object} models.ErrorResponse "URL expired"
// @Router /{code} [get]
func (e *endpoints) UrlRedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	resp, err := e.s.UrlRedirect(r.Context(), code)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}
