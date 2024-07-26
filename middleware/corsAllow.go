package middleware

import (
	"net/http"
)

func MiddlewareCORSAllow() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		var corsOrigin = "*"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				// corsOrigin = cfg["CORS_ORIGIN"]
				w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin, Authorization")
				w.WriteHeader(http.StatusOK)
			}
			next.ServeHTTP(w, r)
		})
	}
}
