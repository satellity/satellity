package views

import (
	"net/http"
	"satellity/internal/models"
)

// VerificationView is the response body of EmailVerification
type VerificationView struct {
	Type           string `json:"type"`
	VerificationID string `json:"verification_id"`
}

func buildVerification(verification *models.EmailVerification) VerificationView {
	return VerificationView{
		Type:           "email_verification",
		VerificationID: verification.VerificationID,
	}
}

// RenderVerification response an email_verification
func RenderVerification(w http.ResponseWriter, r *http.Request, verification *models.EmailVerification) {
	RenderResponse(w, r, buildVerification(verification))
}
