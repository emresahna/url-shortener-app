package ipaddr

import (
	"context"
	"net/http"

	"github.com/EmreSahna/url-shortener-app/internal/models"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if ip == "" {
			http.Error(w, "Ip address is required", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.IpKey, ip)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
