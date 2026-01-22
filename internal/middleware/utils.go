package middleware

import "net/http"

// ToMiddleware преобразует func(http.HandlerFunc) http.HandlerFunc
// в func(http.Handler) http.Handler
func ToMiddleware(fn func(http.HandlerFunc) http.HandlerFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fn(handler.ServeHTTP)(w, r)
		})
	}
}
