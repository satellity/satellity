package middleware

import (
	"net/http"

	"github.com/godiscourse/godiscourse/api/session"
	"github.com/unrolled/render"
)

// Context put database and request in r.Context
func Context(handler http.Handler, r *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := session.WithRender(req.Context(), r)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
