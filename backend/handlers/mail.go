package handlers

import (
	"context"
	"fmt"
	"log"
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

	// Generate random verification code
	resetCode, err := auth.GenerateRandomString(6)
	if err != nil {
		log.Printf("[SendVerifyEmail] Error generating verification code for user %s: %v", user.Email, err)
		return types.VerificationRow{}, fmt.Errorf("failed to generate verification code: %w", err)
	}

	mailData := &mail.MailData{
		Username: user.Username,
		Code:     resetCode,
	}

	// Send verification email
	email := env.MS.NewMail(from, to, "", mailType, mailData)
	err = env.MS.SendMail(context.Background(), email)
	if err != nil {
		log.Printf("[SendVerifyEmail] Error sending verification email to %s: %v", user.Email, err)
		return types.VerificationRow{}, fmt.Errorf("failed to send verification email: %w", err)
	}

	// Hash the verification code before storing
	hashedCode, err := auth.GetHash(resetCode)
	if err != nil {
		log.Printf("[SendVerifyEmail] Error hashing verification code for user %s: %v", user.Email, err)
		return types.VerificationRow{}, fmt.Errorf("failed to hash verification code: %w", err)
	}

	log.Printf("[SendVerifyEmail] Successfully sent verification email to %s", user.Email)

	return types.VerificationRow{
		Email:     user.Email,
		Code:      hashedCode,
		ExpiresAt: time.Now().Add(time.Hour),
		Type:      mailType,
	}, nil
}

func (env *Env) SendPasswordResetEmail(user types.UserRow) (types.VerificationRow, error) {
	from := os.Getenv("PLAYABILITY_EMAIL")
	to := []string{user.Email}
	mailType := mail.PassReset

	// Generate random password reset code
	resetCode, err := auth.GenerateRandomString(6)
	if err != nil {
		log.Printf("[SendPasswordResetEmail] Error generating reset code for user %s: %v", user.Email, err)
		return types.VerificationRow{}, fmt.Errorf("failed to generate reset code: %w", err)
	}

	mailData := &mail.MailData{
		Username: user.Username,
		Code:     resetCode,
	}

	// Send password reset email
	email := env.MS.NewMail(from, to, "", mailType, mailData)
	err = env.MS.SendMail(context.Background(), email)
	if err != nil {
		log.Printf("[SendPasswordResetEmail] Error sending reset email to %s: %v", user.Email, err)
		return types.VerificationRow{}, fmt.Errorf("failed to send password reset email: %w", err)
	}

	// Hash the reset code before storing
	hashedCode, err := auth.GetHash(resetCode)
	if err != nil {
		log.Printf("[SendPasswordResetEmail] Error hashing reset code for user %s: %v", user.Email, err)
		return types.VerificationRow{}, fmt.Errorf("failed to hash reset code: %w", err)
	}

	log.Printf("[SendPasswordResetEmail] Successfully sent password reset email to %s", user.Email)

	return types.VerificationRow{
		Email:     user.Email,
		Code:      hashedCode,
		ExpiresAt: time.Now().Add(time.Hour),
		Type:      mailType,
	}, nil
}

