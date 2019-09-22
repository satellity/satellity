package middlewares

import (
	"context"
	"net/http"
	"regexp"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"strings"
)

var whitelist = [][2]string{
	{"GET", "^/_hc$"},
	{"GET", "^/categories"},
	{"GET", "^/topics"},
	{"GET", "^/users"},
	{"GET", "^/groups$"},
	{"GET", "^/groups/[a-z0-9-]+$"},
	{"POST", "^/oauth"},
	{"POST", "^/sessions"},
	{"POST", "^/email_verifications"},
}

var userWhitelist = [][2]string{
	{"GET", "^/me"},
	{"POST", "^/comments"},
	{"DELETE", "^/comments"},
	{"POST", "^/topics"},
	{"POST", "^/me"},
	{"POST", "^/groups"},
	{"GET", "^/groups"},
	{"GET", "^/user"},
}

type contextValueKey int

const keyCurrentUser contextValueKey = 1000

// CurrentUser read the user from r.Context
func CurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(keyCurrentUser).(*models.User)
	return user
}

// Authenticate handle routes by user's role
func Authenticate(database *durable.Database, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			handleUnauthorized(handler, w, r)
			return
		}
		mcontext := models.WrapContext(r.Context(), database)
		user, err := models.AuthenticateUser(mcontext, header[7:])
		if err != nil {
			views.RenderErrorResponse(w, r, err)
			return
		}
		if user == nil {
			handleUnauthorized(handler, w, r)
			return
		}
		ctx := context.WithValue(r.Context(), keyCurrentUser, user)
		if user.Role() != "admin" {
			handleUserRouters(handler, w, r.WithContext(ctx))
			return
		}
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleUnauthorized(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	for _, pp := range whitelist {
		if pp[0] != r.Method {
			continue
		}
		if matched, _ := regexp.MatchString(pp[1], strings.ToLower(r.URL.Path)); matched {
			handler.ServeHTTP(w, r)
			return
		}
	}

	views.RenderErrorResponse(w, r, session.AuthorizationError(r.Context()))
}

func handleUserRouters(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	for _, pp := range userWhitelist {
		if pp[0] != r.Method {
			continue
		}
		if matched, _ := regexp.MatchString(pp[1], strings.ToLower(r.URL.Path)); matched {
			handler.ServeHTTP(w, r)
			return
		}
	}

	handleUnauthorized(handler, w, r)
}
