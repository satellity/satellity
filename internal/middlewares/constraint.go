package middlewares

import (
	"net/http"
	"satellity/internal/views"
)

// Constraint process OPTIONS request.
func Constraint(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST,DELETE")
			w.Header().Set("Access-Control-Max-Age", "86400")
		}
		if r.Method == "OPTIONS" {
			views.RenderBlankResponse(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
