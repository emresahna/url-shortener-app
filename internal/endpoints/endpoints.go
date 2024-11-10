package endpoints

import (
	"encoding/json"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Endpoints interface {
	SignupUserHandler(w http.ResponseWriter, r *http.Request)
	LoginUserHandler(w http.ResponseWriter, r *http.Request)
	UrlShortenerHandler(w http.ResponseWriter, r *http.Request)
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
		e.encodeResponse(w, err)
		return
	}

	e.encodeResponse(w, resp)
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
		e.encodeResponse(w, err)
		return
	}

	e.encodeResponse(w, resp)
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
		e.encodeResponse(w, err)
		return
	}

	e.encodeResponse(w, resp)
	return
}

func (e *endpoints) RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	resp, err := e.s.RedirectUrl(r.Context(), code)
	if err != nil {
		e.encodeResponse(w, err)
		return
	}

	w.Header().Set("Location", resp)
	w.WriteHeader(301)
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

func (e *endpoints) encodeResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
