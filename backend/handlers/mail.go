package handlers

import (
	"context"
	"os"
	"playability/auth"
	"playability/pkg/mail"
	"playability/types"
	"time"
)

func (env *Env) SendVerifyEmail(user types.UserRow) (types.VerificationRow, error) {
	from := os.Getenv("PLAYABILITY_EMAIL")
	to := []string{user.Email}
	mailType := mail.MailConfirmation
	resetCode, err := auth.GenerateRandomString(6)
	if err != nil {
		//todo error handling
	}
	mailData := &mail.MailData{
		Username: user.Username,
		Code:     resetCode,
	}
	email := env.MS.NewMail(from, to, "", mailType, mailData)
	err = env.MS.SendMail(context.Background(), email)
	if err != nil {
		//todo error handling
	}
	hashedCode, err := auth.GetHash(resetCode)
	if err != nil {
		//todo error handling
	}
	return types.VerificationRow{
		Email:     user.Email,
		Code:      hashedCode,
		ExpiresAt: time.Now().Add(time.Hour),
		Type:      mailType,
	}, nil
}

