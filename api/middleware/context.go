package middleware

import (
	"database/sql"
	"net/http"

	"github.com/godiscourse/godiscourse/api/session"
	"github.com/unrolled/render"
)

// Context put database and request in r.Context
func Context(handler http.Handler, db *sql.DB, r *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := session.WithDatabase(req.Context(), db)
		ctx = session.WithRender(ctx, r)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
