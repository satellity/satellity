package middleware

import (
	"net/http"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/session"
	"github.com/unrolled/render"
)

func Context(handler http.Handler, db *pg.DB, rend *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := session.WithDatabase(r.Context(), db)
		ctx = session.WithRender(ctx, rend)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
