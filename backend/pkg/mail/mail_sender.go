package mail

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/mailersend/mailersend-go"
)

// Implement MailService for Mail Sender
type MSMailService struct {
	service *mailersend.Mailersend
	logger  hclog.Logger
}

func NewMSMailService(logger hclog.Logger, apiKey string) *MSMailService {
	ms := mailersend.NewMailersend(apiKey)
	return &MSMailService{ms, logger}
}

// createMail takes in a mail request and constructs a Mail Sender mail type
func (ms *MSMailService) createMail(mail *Mail) (*mailersend.Message, error) {
	subject := ""
	templateID := ""
	tags := []string{}
	switch mail.mtype {
	case MailConfirmation:
		subject = "Confirm your email for Playability"
		templateID = os.Getenv("CONFIRMATION_TEMPLATE_ID")
		if templateID == "" {
			ms.logger.Error("CONFIRMATION_TEMPLATE_ID environment variable is not set")
			return nil, ErrMissingTemplateID
		}
		tags = []string{"email_confirmation"}
	case PassReset:
		subject = "Reset your Playability account password"
		templateID = os.Getenv("PASS_RESET_TEMPLATE_ID")
		if templateID == "" {
			ms.logger.Error("PASS_RESET_TEMPLATE_ID environment variable is not set")
			return nil, ErrMissingTemplateID
		}
		tags = []string{"password_reset"}
	default:
		// Use custom subject for non-template emails
		subject = mail.subject
	}
	from := mailersend.From{
		Email: mail.from,
	}
	to := []mailersend.Recipient{}
	for _, recipient := range mail.to {
		to = append(to, mailersend.Recipient{
			Email: recipient,
		})
	}
	personalization := []mailersend.Personalization{
		{
			Data: map[string]interface{}{
				"code":     mail.data.Code,
				"username": mail.data.Username,
			},
		},
	}
	message := ms.service.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(to)
	message.SetSubject(subject)
	message.SetTemplateID(templateID)
	message.SetPersonalization(personalization)
	message.SetTags(tags)

	return message, nil
}
func (ms *MSMailService) SendMail(ctx context.Context, mail *Mail) error {
	msmail, err := ms.createMail(mail)
	if err != nil {
		ms.logger.Error("failed to create mail", "error", err, "mail_type", mail.mtype)
		return err
	}

	res, err := ms.service.Email.Send(ctx, msmail)
	if err != nil {
		ms.logger.Error("failed to send mail", "error", err, "mail_type", mail.mtype, "recipients", mail.to)
		return err
	}

	ms.logger.Info("mail sent successfully", "mail_type", mail.mtype, "recipients", mail.to, "message_id", res.Header.Get("X-Message-Id"))
	return nil
}
func (ms *MSMailService) NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail {
	return &Mail{
		from:    from,
		to:      to,
		subject: subject,
		mtype:   mailType,
		data:    data,
	}
}
