package middleware

import (
	"satellity/internal/durable"
	"satellity/internal/session"
	"net/http"
)

// Logger put logger in r.Context
func Logger(handler http.Handler, logger *durable.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := session.WithLogger(r.Context(), logger)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
