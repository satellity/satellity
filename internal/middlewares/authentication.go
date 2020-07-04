package middlewares

import (
	"context"
	"net/http"
	"regexp"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"strings"
)

var whitelist = [][2]string{
	{"GET", "^/api/_hc$"},
	{"GET", "^/api/products"},
	{"GET", "^/api/categories"},
	{"GET", "^/api/client"},
	{"GET", "^/api/topics"},
	{"GET", "^/api/users"},
	{"POST", "^/api/oauth"},
	{"POST", "^/api/sessions"},
	{"POST", "^/api/email_verifications"},
}

var userWhitelist = [][2]string{
	{"GET", "^/api/me"},
	{"POST", "^/api/comments"},
	{"DELETE", "^/api/comments"},
	{"POST", "^/api/topics"},
	{"POST", "^/api/me"},
	{"GET", "^/api/user"},
}

type contextValueKey int

const keyCurrentUser contextValueKey = 1000

// CurrentUser read the user from r.Context
func CurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(keyCurrentUser).(*models.User)
	return user
}

// Authenticate handle routes by user's role
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
		if user.GetRole() != models.UserRoleAdmin {
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
