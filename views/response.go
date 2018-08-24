package views

import (
	"net/http"

	"github.com/godiscourse/godiscourse/session"
)

type ResponseView struct {
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"error,omitempty"`
}

func RenderResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	session.Render(r.Context()).JSON(w, http.StatusOK, ResponseView{Data: data})
}

func RenderErrorResponse(w http.ResponseView, r *http.Request, err error) {
	sessionError, ok := err.(session.Error)
	if !ok {
		sessionError = session.ServerError(r.Context(), err)
	}
	session.Render(r.Context()).JSON(w, sessionError.Status, ResponseView{Error: sessionError})
}

func RenderBlankResponse(w http.ResponseWriter, r *http.Request) {
	session.Render(r.Context()).JSON(w, http.StatusOK, map[string]string{})
}
