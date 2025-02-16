package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/service"
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
}

type endpoints struct {
	s service.Service
}

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

func (e *endpoints) UserMeHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := e.s.UserMe(r.Context())
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}

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

func (e *endpoints) UrlRemoveHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := e.s.UrlRemove(r.Context(), id)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}

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

func (e *endpoints) UrlRedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	resp, err := e.s.UrlRedirect(r.Context(), code)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
}

func NewEndpoints(s service.Service) Endpoints {
	return &endpoints{
		s: s,
	}
}

func (e *endpoints) decodeRequest(body io.ReadCloser, req interface{}) (err error) {
	if err := json.NewDecoder(body).Decode(req); err != nil {
		return err
	}
	return nil
}

func (e *endpoints) encodeResponse(w http.ResponseWriter, res interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		e.encodeError(w, fmt.Errorf("failed to encode response: %w", err))
		return
	}
}

func (e *endpoints) encodeError(w http.ResponseWriter, err error) {
	var apiErr *models.Error

	if ok := errors.As(err, &apiErr); !ok {
		apiErr = models.InternalServerErr()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.StatusCode)

	if encodeErr := json.NewEncoder(w).Encode(&apiErr); encodeErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
