package handlers

import (
	"playability/db"
	"playability/pkg/ai"
	"playability/pkg/fetch"
	"playability/pkg/mail"
)

type Env struct {
	DB db.DatabaseModel
	MS mail.MailService
	AI ai.AIService
	FS fetch.FetchService
}
