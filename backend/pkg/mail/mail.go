package mail

import (
	"context"
	"errors"
)

var (
	ErrMissingTemplateID = errors.New("mail template ID not configured")
)

type MailData struct {
	Username string
	Code     string
}
type MailType int

const (
	MailConfirmation MailType = iota + 1
	PassReset
)

type MailUser struct {
	name  string
	email string
}

type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	mtype   MailType
	data    *MailData
}
type MailService interface {
	SendMail(ctx context.Context, mail *Mail) error
	NewMail(from string, to []string, subject string, mtype MailType, data *MailData) *Mail
}
