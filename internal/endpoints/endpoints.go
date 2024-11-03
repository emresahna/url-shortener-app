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
	GetUserHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
	UrlShortenerHandler(w http.ResponseWriter, r *http.Request)
}

type endpoints struct {
	s service.Service
}

func NewEndpoints(s service.Service) Endpoints {
	return &endpoints{
		s: s,
	}
}

func (e *endpoints) DecodeRequest(body io.ReadCloser, req interface{}) (err error) {
	if err := json.NewDecoder(body).Decode(req); err != nil {
		return err
	}
	return nil
}

func (e *endpoints) EncodeResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (e *endpoints) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var userId string
	userId = chi.URLParam(r, "userID")

	resp, err := e.s.GetUser(r.Context(), userId)
	if err != nil {
		e.EncodeResponse(w, err)
		return
	}

	e.EncodeResponse(w, resp)
	return
}

func (e *endpoints) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	err := e.DecodeRequest(r.Body, &req)
	if err != nil {
		return
	}

	resp, err := e.s.CreateUser(r.Context(), req)
	if err != nil {
		e.EncodeResponse(w, err)
		return
	}

	e.EncodeResponse(w, resp)
	return
}

func (e *endpoints) UrlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenURLRequest
	err := e.DecodeRequest(r.Body, &req)
	if err != nil {
		return
	}

	resp, err := e.s.ShortenURL(r.Context(), req)
	if err != nil {
		e.EncodeResponse(w, err)
		return
	}

	e.EncodeResponse(w, resp)
	return
}
