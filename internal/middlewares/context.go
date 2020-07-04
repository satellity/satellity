package middlewares

import (
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/session"

	"github.com/unrolled/render"
)

// Context put database and request in r.Context
func Context(handler http.Handler, d *durable.Database, r *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := session.WithDatabase(req.Context(), d)
		ctx = session.WithRender(ctx, r)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
