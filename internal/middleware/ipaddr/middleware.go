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

		tempIp := "127.0.0.1"

		ctx := context.WithValue(r.Context(), "ip", tempIp)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
