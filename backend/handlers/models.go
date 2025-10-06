package handlers

import (
	"playability/db"
	"playability/pkg/mail"
)

type Env struct {
	DB db.DatabaseModel
	MS mail.MailService
}
