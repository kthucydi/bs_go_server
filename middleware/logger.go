package middleware

import (
	"net/http"
)

func MiddlewareSmartLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logger.Printf("InnerLog:Request path: %s, method: %s", r.URL.Path, r.Method)
			next.ServeHTTP(w, r)
		})
	}
}
