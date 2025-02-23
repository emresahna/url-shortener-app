package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/emresahna/url-shortener-app/internal/models"
)

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
