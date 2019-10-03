package clouds

import (
	"context"
	"fmt"
	"satellity/internal/configs"
	"time"

	mailgun "github.com/mailgun/mailgun-go/v3"
)

// SendVerificationEmail send an verification email
func SendVerificationEmail(ctx context.Context, purpose, recipient, code string) error {
	v := configs.AppConfig.Email.Verification
	title := v.Title
	if purpose == "PASSWORD" {
		title = v.Reset
	}
	return sendEmail(ctx, title, fmt.Sprintf(v.Body, code), recipient)
}

func sendEmail(ctx context.Context, subject, body, recipient string) error {
	config := configs.AppConfig
	if config.Environment == "test" {
		return nil
	}
	mg := mailgun.NewMailgun(config.Mailgun.Domain, config.Mailgun.Key)
	sender := config.Mailgun.Sender

	message := mg.NewMessage(sender, subject, "", recipient)
	message.SetHtml(body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	return err
}
