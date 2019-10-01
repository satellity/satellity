package clouds

import (
	"context"
	"fmt"
	"satellity/internal/configs"
	"time"

	mailgun "github.com/mailgun/mailgun-go"
)

// SendVerificationEmail send an verification email
func SendVerificationEmail(ctx context.Context, recipient, code string) error {
	v := configs.AppConfig.Email.Verification
	return sendEmail(ctx, v.Title, fmt.Sprintf(v.Body, code), recipient)
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
