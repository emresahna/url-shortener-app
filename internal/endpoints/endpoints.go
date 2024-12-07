package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Endpoints interface {
	SignupUserHandler(w http.ResponseWriter, r *http.Request)
	LoginUserHandler(w http.ResponseWriter, r *http.Request)
	RefreshHandler(w http.ResponseWriter, r *http.Request)
	MeHandler(w http.ResponseWriter, r *http.Request)
	UrlShortenerHandler(w http.ResponseWriter, r *http.Request)
	RemoveUrlHandler(w http.ResponseWriter, r *http.Request)
	LimitedUrlShortenerHandler(w http.ResponseWriter, r *http.Request)
	RedirectUrlHandler(w http.ResponseWriter, r *http.Request)
}

type endpoints struct {
	s service.Service
}

func (e *endpoints) SignupUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SignupUserRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		return
	}

	resp, err := e.s.SignupUser(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 201)
	return
}

func (e *endpoints) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		return
	}

	resp, err := e.s.LoginUser(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
	return
}

func (e *endpoints) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		return
	}

	resp, err := e.s.Refresh(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
	return
}

func (e *endpoints) MeHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := e.s.Me(r.Context())
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
	return
}

func (e *endpoints) UrlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenURLRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		return
	}

	resp, err := e.s.ShortenURL(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 201)
	return
}

func (e *endpoints) RemoveUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := e.s.RemoveUrl(r.Context(), id)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
	return
}

func (e *endpoints) LimitedUrlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenURLRequest
	err := e.decodeRequest(r.Body, &req)
	if err != nil {
		return
	}

	resp, err := e.s.LimitedShortenURL(r.Context(), req)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 201)
	return
}

func (e *endpoints) RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	resp, err := e.s.RedirectUrl(r.Context(), code)
	if err != nil {
		e.encodeError(w, err)
		return
	}

	e.encodeResponse(w, resp, 200)
	return
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
	json.NewEncoder(w).Encode(res)
}

func (e *endpoints) encodeError(w http.ResponseWriter, err error) {
	var apiErr *models.Error

	if ok := errors.As(err, &apiErr); !ok {
		apiErr = models.InternalServerErr()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiErr.StatusCode)
		json.NewEncoder(w).Encode(apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.StatusCode)
	json.NewEncoder(w).Encode(apiErr)
}
