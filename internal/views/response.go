package views

import (
	"net/http"
	"satellity/internal/session"
)

// ResponseView is the struct of response
type ResponseView struct {
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"error,omitempty"`
}

// RenderResponse respond a special data, e.g.: Topic, Category etc.
func RenderResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	session.Render(r.Context()).JSON(w, http.StatusOK, ResponseView{Data: data})
}

// RenderErrorResponse respond an error response
func RenderErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sessionError, ok := err.(session.Error)
	if !ok {
		sessionError = session.ServerError(r.Context(), err)
	}
	if sessionError.Code == 10001 {
		sessionError.Code = 500
	}
	session.Render(r.Context()).JSON(w, sessionError.Status, ResponseView{Error: sessionError})
}

// RenderBlankResponse respond a blank response
func RenderBlankResponse(w http.ResponseWriter, r *http.Request) {
	session.Render(r.Context()).JSON(w, http.StatusOK, map[string]string{})
}
