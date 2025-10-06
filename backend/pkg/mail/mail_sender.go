package mail

import (
	"context"
	"fmt"
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
	logger.Info("Initializing MailSend service")
	if apiKey == "" {
		logger.Warn("API key is empty - mail sending will fail")
	} else {
		logger.Info("API key provided", "key_length", len(apiKey))
	}
	ms := mailersend.NewMailersend(apiKey)
	logger.Info("MailSend service initialized successfully")
	return &MSMailService{ms, logger}
}

// createMail takes in a mail request and constructs a Mail Sender mail type
func (ms *MSMailService) createMail(mail *Mail) (*mailersend.Message, error) {
	ms.logger.Info("Creating email message", "mail_type", mail.mtype, "recipients", mail.to, "from", mail.from)

	subject := ""
	templateID := ""
	tags := []string{}
	switch mail.mtype {
	case MailConfirmation:
		ms.logger.Debug("Processing email confirmation mail type")
		subject = "Confirm your email for Playability"
		templateID = os.Getenv("CONFIRMATION_TEMPLATE_ID")
		ms.logger.Debug("Template ID from env", "template_id", templateID, "is_empty", templateID == "")
		if templateID == "" {
			ms.logger.Error("CONFIRMATION_TEMPLATE_ID environment variable is not set")
			return nil, ErrMissingTemplateID
		}
		tags = []string{"email_confirmation"}
	case PassReset:
		ms.logger.Debug("Processing password reset mail type")
		subject = "Reset your Playability account password"
		templateID = os.Getenv("PASS_RESET_TEMPLATE_ID")
		ms.logger.Debug("Template ID from env", "template_id", templateID, "is_empty", templateID == "")
		if templateID == "" {
			ms.logger.Error("PASS_RESET_TEMPLATE_ID environment variable is not set")
			return nil, ErrMissingTemplateID
		}
		tags = []string{"password_reset"}
	default:
		ms.logger.Warn("Unknown mail type, using custom subject", "mail_type", mail.mtype)
		// Use custom subject for non-template emails
		subject = mail.subject
	}
	from := mailersend.From{
		Email: mail.from,
	}
	ms.logger.Debug("Setting from address", "from", mail.from)

	to := []mailersend.Recipient{}
	for _, recipient := range mail.to {
		to = append(to, mailersend.Recipient{
			Email: recipient,
		})
	}
	ms.logger.Debug("Setting recipients", "count", len(to), "recipients", mail.to)

	personalization := []mailersend.Personalization{
		{
			Data: map[string]interface{}{
				"code":     mail.data.Code,
				"username": mail.data.Username,
			},
		},
	}
	ms.logger.Debug("Setting personalization", "username", mail.data.Username, "has_code", mail.data.Code != "")

	ms.logger.Debug("Building message", "subject", subject, "template_id", templateID, "tags", tags)
	message := ms.service.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(to)
	message.SetSubject(subject)
	message.SetTemplateID(templateID)
	message.SetPersonalization(personalization)
	message.SetTags(tags)

	ms.logger.Info("Email message created successfully", "subject", subject, "template_id", templateID, "recipient_count", len(to))
	return message, nil
}
func (ms *MSMailService) SendMail(ctx context.Context, mail *Mail) error {
	ms.logger.Info("SendMail called", "mail_type", mail.mtype, "recipients", mail.to)

	msmail, err := ms.createMail(mail)
	if err != nil {
		ms.logger.Error("failed to create mail", "error", err, "mail_type", mail.mtype)
		return err
	}

	ms.logger.Info("Attempting to send email via MailerSend API", "mail_type", mail.mtype, "recipients", mail.to)
	res, err := ms.service.Email.Send(ctx, msmail)
	if err != nil {
		ms.logger.Error("failed to send mail via MailerSend API", "error", err, "error_type", fmt.Sprintf("%T", err), "mail_type", mail.mtype, "recipients", mail.to)
		return err
	}

	messageID := res.Header.Get("X-Message-Id")
	ms.logger.Info("mail sent successfully via MailerSend",
		"mail_type", mail.mtype,
		"recipients", mail.to,
		"message_id", messageID,
		"status_code", res.StatusCode)
	return nil
}
func (ms *MSMailService) NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail {
	ms.logger.Debug("NewMail called",
		"from", from,
		"to", to,
		"subject", subject,
		"mail_type", mailType,
		"username", data.Username,
		"has_code", data.Code != "")

	mail := &Mail{
		from:    from,
		to:      to,
		subject: subject,
		mtype:   mailType,
		data:    data,
	}

	ms.logger.Debug("Mail object created successfully", "mail_type", mailType, "recipients", to)
	return mail
}
