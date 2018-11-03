package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/views"
)

type userRequest struct {
	Code          string `json:"code"`
	SessionSecret string `json:"session_secret"`
}

type userImpl struct{}

func registerUser(router *httptreemux.TreeMux) {
	impl := &userImpl{}
	router.POST("/oauth/:provider", impl.oauth)
}

func (impl *userImpl) oauth(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body userRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	user, err := models.CreateGithubUser(r.Context(), body.Code, body.SessionSecret)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderAccount(w, r, user)
}
