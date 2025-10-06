package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"playability/auth"
	"playability/pkg/mail"
	"strings"
	"time"
)

// VerifyEmailRequest represents the request body for email verification
type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// PasswordResetRequest represents the request body for password reset
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents the request body for resetting password with code
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

// PostVerifyEmail handles email verification with a code
func (env *Env) PostVerifyEmail(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PostVerifyEmail] Received verification request from %s", r.RemoteAddr)

	var req VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[PostVerifyEmail] Invalid JSON in request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Code == "" {
		log.Printf("[PostVerifyEmail] Missing required fields")
		http.Error(w, "email and code are required", http.StatusBadRequest)
		return
	}

	// Normalize email
	req.Email = strings.ToLower(req.Email)

	log.Printf("[PostVerifyEmail] Verifying email: %s", req.Email)

	// Get verification from database
	verification, err := env.DB.GetVerification(req.Email, int(mail.MailConfirmation))
	if err != nil {
		log.Printf("[PostVerifyEmail] Error getting verification: %v", err)
		http.Error(w, "invalid or expired verification code", http.StatusBadRequest)
		return
	}

	// Check if verification has expired
	if time.Now().After(verification.ExpiresAt) {
		log.Printf("[PostVerifyEmail] Verification code expired for: %s", req.Email)
		// Clean up expired verification
		env.DB.DeleteVerification(req.Email, int(mail.MailConfirmation))
		http.Error(w, "verification code has expired", http.StatusBadRequest)
		return
	}

	// Verify the code
	err = auth.CheckPassword(req.Code, verification.Code)
	if err != nil {
		log.Printf("[PostVerifyEmail] Invalid verification code for: %s", req.Email)
		http.Error(w, "invalid verification code", http.StatusBadRequest)
		return
	}

	// Mark user as verified
	err = env.DB.MarkUserAsVerified(req.Email)
	if err != nil {
		log.Printf("[PostVerifyEmail] Error marking user as verified: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Delete the verification code
	err = env.DB.DeleteVerification(req.Email, int(mail.MailConfirmation))
	if err != nil {
		log.Printf("[PostVerifyEmail] Error deleting verification: %v", err)
		// Don't return error as user is already verified
	}

	log.Printf("[PostVerifyEmail] Successfully verified email: %s", req.Email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "email verified successfully"})
}

// PostRequestPasswordReset handles password reset email requests
func (env *Env) PostRequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PostRequestPasswordReset] Received password reset request from %s", r.RemoteAddr)

	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[PostRequestPasswordReset] Invalid JSON in request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" {
		log.Printf("[PostRequestPasswordReset] Missing email")
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	// Normalize email
	req.Email = strings.ToLower(req.Email)

	log.Printf("[PostRequestPasswordReset] Password reset requested for: %s", req.Email)

	// Get user by email
	user, err := env.DB.GetUserByEmail(req.Email)
	if err != nil {
		// For security reasons, don't reveal if user exists or not
		log.Printf("[PostRequestPasswordReset] User not found: %s", req.Email)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "if the email exists, a password reset link has been sent"})
		return
	}

	// Send password reset email
	verification, err := env.SendPasswordResetEmail(user)
	if err != nil {
		log.Printf("[PostRequestPasswordReset] Error sending password reset email: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Store verification in database
	err = env.DB.InsertVerification(verification)
	if err != nil {
		log.Printf("[PostRequestPasswordReset] Error storing verification: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[PostRequestPasswordReset] Password reset email sent to: %s", req.Email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "if the email exists, a password reset link has been sent"})
}

// PostResetPassword handles password reset with verification code
func (env *Env) PostResetPassword(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PostResetPassword] Received password reset from %s", r.RemoteAddr)

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[PostResetPassword] Invalid JSON in request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Code == "" || req.NewPassword == "" {
		log.Printf("[PostResetPassword] Missing required fields")
		http.Error(w, "email, code, and new_password are required", http.StatusBadRequest)
		return
	}

	// Normalize email
	req.Email = strings.ToLower(req.Email)

	log.Printf("[PostResetPassword] Resetting password for: %s", req.Email)

	// Get verification from database
	verification, err := env.DB.GetVerification(req.Email, int(mail.PassReset))
	if err != nil {
		log.Printf("[PostResetPassword] Error getting verification: %v", err)
		http.Error(w, "invalid or expired verification code", http.StatusBadRequest)
		return
	}

	// Check if verification has expired
	if time.Now().After(verification.ExpiresAt) {
		log.Printf("[PostResetPassword] Verification code expired for: %s", req.Email)
		// Clean up expired verification
		env.DB.DeleteVerification(req.Email, int(mail.PassReset))
		http.Error(w, "verification code has expired", http.StatusBadRequest)
		return
	}

	// Verify the code
	err = auth.CheckPassword(req.Code, verification.Code)
	if err != nil {
		log.Printf("[PostResetPassword] Invalid verification code for: %s", req.Email)
		http.Error(w, "invalid verification code", http.StatusBadRequest)
		return
	}

	// Hash the new password
	newHash, err := auth.GetHash(req.NewPassword)
	if err != nil {
		log.Printf("[PostResetPassword] Error hashing password: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Update user password
	err = env.DB.UpdateUserPassword(req.Email, newHash)
	if err != nil {
		log.Printf("[PostResetPassword] Error updating password: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Delete the verification code
	err = env.DB.DeleteVerification(req.Email, int(mail.PassReset))
	if err != nil {
		log.Printf("[PostResetPassword] Error deleting verification: %v", err)
		// Don't return error as password is already reset
	}

	log.Printf("[PostResetPassword] Successfully reset password for: %s", req.Email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "password reset successfully"})
}

// PostResendVerification handles resending verification emails
func (env *Env) PostResendVerification(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PostResendVerification] Received resend verification request from %s", r.RemoteAddr)

	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[PostResendVerification] Invalid JSON in request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" {
		log.Printf("[PostResendVerification] Missing email")
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	// Normalize email
	req.Email = strings.ToLower(req.Email)

	log.Printf("[PostResendVerification] Resend verification requested for: %s", req.Email)

	// Get user by email
	user, err := env.DB.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("[PostResendVerification] User not found: %s", req.Email)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Check if user is already verified
	if user.Verified {
		log.Printf("[PostResendVerification] User already verified: %s", req.Email)
		http.Error(w, "email is already verified", http.StatusBadRequest)
		return
	}

	// Send verification email
	verification, err := env.SendVerifyEmail(user)
	if err != nil {
		log.Printf("[PostResendVerification] Error sending verification email: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Store verification in database
	err = env.DB.InsertVerification(verification)
	if err != nil {
		log.Printf("[PostResendVerification] Error storing verification: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[PostResendVerification] Verification email resent to: %s", req.Email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "verification email sent"})
}
