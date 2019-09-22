package controllers

import (
	"encoding/json"
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type verificationImpl struct {
	database *durable.Database
}

type verificationRequest struct {
	Recaptcha     string `json:"recaptcha"`
	Email         string `json:"email"`
	Code          string `json:"code"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	SessionSecret string `json:"session_secret"`
}

func registerVerification(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &verificationImpl{database: database}

	router.POST("/email_verifications", impl.create)
	router.POST("/email_verifications/:id", impl.verify)
}

func (impl *verificationImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body verificationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if verification, err := models.CreateEmailVerification(mctx, body.Email, body.Recaptcha); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderVerification(w, r, verification)
	}
}

func (impl *verificationImpl) verify(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body verificationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	user, err := models.VerifyEmailVerification(mctx, params["id"], body.Code, body.Username, body.Password, body.SessionSecret)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAccount(w, r, user)
	}
}
