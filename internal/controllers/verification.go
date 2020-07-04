package controllers

import (
	"encoding/json"
	"net/http"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

const (
	purposeUser     = "USER"
	purposePassword = "PASSWORD"
)

type verificationImpl struct{}

type verificationRequest struct {
	Recaptcha     string `json:"recaptcha"`
	Email         string `json:"email"`
	Code          string `json:"code"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Purpose       string `json:"purpose"`
	SessionSecret string `json:"session_secret"`
}

func registerVerification(router *httptreemux.Group) {
	impl := &verificationImpl{}

	router.POST("/email_verifications", impl.create)
	router.POST("/email_verifications/:id", impl.verify)
}

func (impl *verificationImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body verificationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if verification, err := models.CreateEmailVerification(r.Context(), body.Purpose, body.Email, body.Recaptcha); err != nil {
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

	switch body.Purpose {
	case purposeUser:
		user, err := models.VerifyEmailVerification(r.Context(), params["id"], body.Code, body.Username, body.Password, body.SessionSecret)
		if err != nil {
			views.RenderErrorResponse(w, r, err)
		} else {
			views.RenderAccount(w, r, user)
		}
	case purposePassword:
		err := models.Reset(r.Context(), params["id"], body.Code, body.Password)
		if err != nil {
			views.RenderErrorResponse(w, r, err)
		} else {
			views.RenderBlankResponse(w, r)
		}
	default:
		views.RenderErrorResponse(w, r, session.BadDataError(r.Context()))
	}
}
