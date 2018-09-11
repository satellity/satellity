package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/views"
)

var whitelist = map[string]string{
	"POST": "^/users$",
}

type contextValueKey int

const keyCurrentUser contextValueKey = 1000

func CurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(keyCurrentUser).(*models.User)
	return user
}

func Authenticate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			handleUnauthorized(handler, w, r)
			return
		}
		user, err := models.AuthenticateUser(r.Context(), header[7:])
		if err != nil {
			views.RenderErrorResponse(w, r, err)
			return
		}
		if user == nil {
			handleUnauthorized(handler, w, r)
			return
		}
		ctx := context.WithValue(r.Context(), keyCurrentUser, user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleUnauthorized(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	for k, v := range whitelist {
		if k != r.Method {
			continue
		}
		if matched, _ := regexp.MatchString(v, strings.ToLower(r.URL.Path)); matched {
			handler.ServeHTTP(w, r)
			return
		}
	}

	views.RenderErrorResponse(w, r, session.AuthorizationError(r.Context()))
}
