package ipaddr

import (
	"context"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if ip == "" {
			http.Error(w, "Ip address is required", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "ip", ip)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
