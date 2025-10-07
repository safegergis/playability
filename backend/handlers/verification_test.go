package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"playability/auth"
	"playability/db"
	"playability/pkg/mail"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// MockMailService implements mail.MailService for testing
type MockMailService struct {
	ShouldFail bool
	SentMails  []*mail.Mail
}

func (m *MockMailService) SendMail(ctx context.Context, mailObj *mail.Mail) error {
	if m.ShouldFail {
		return errors.New("mock mail service error")
	}
	m.SentMails = append(m.SentMails, mailObj)
	return nil
}

func (m *MockMailService) NewMail(from string, to []string, subject string, mtype mail.MailType, data *mail.MailData) *mail.Mail {
	return &mail.Mail{}
}

func TestPostVerifyEmail(t *testing.T) {
	t.Run("successful email verification", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		// Generate valid verification code hash
		codeHash, _ := auth.GetHash("123456")

		req := VerifyEmailRequest{
			Email: "test@example.com",
			Code:  "123456",
		}
		body, _ := json.Marshal(req)

		// Expect GetVerification query
		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnRows(sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
				AddRow("test@example.com", codeHash, time.Now().Add(1*time.Hour), int(mail.MailConfirmation)))

		// Expect MarkUserAsVerified query
		mock.ExpectExec("UPDATE users SET verified = true WHERE email").
			WithArgs("test@example.com").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect DeleteVerification query
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("PostVerifyEmail() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		if response["message"] != "email verified successfully" {
			t.Errorf("PostVerifyEmail() message = %s, want 'email verified successfully'", response["message"])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostVerifyEmail() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "Invalid request body") {
			t.Errorf("PostVerifyEmail() body = %s", w.Body.String())
		}
	})

	t.Run("missing email field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := VerifyEmailRequest{
			Code: "123456",
		}
		body, _ := json.Marshal(req)

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostVerifyEmail() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "email and code are required") {
			t.Errorf("PostVerifyEmail() body = %s", w.Body.String())
		}
	})

	t.Run("missing code field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := VerifyEmailRequest{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(req)

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostVerifyEmail() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("email normalization to lowercase", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		codeHash, _ := auth.GetHash("123456")

		req := VerifyEmailRequest{
			Email: "Test@Example.COM",
			Code:  "123456",
		}
		body, _ := json.Marshal(req)

		// Expect lowercase email
		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnRows(sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
				AddRow("test@example.com", codeHash, time.Now().Add(1*time.Hour), int(mail.MailConfirmation)))

		mock.ExpectExec("UPDATE users SET verified = true WHERE email").
			WithArgs("test@example.com").
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("PostVerifyEmail() status = %d, want %d", w.Code, http.StatusOK)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("expired verification code", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		codeHash, _ := auth.GetHash("123456")

		req := VerifyEmailRequest{
			Email: "test@example.com",
			Code:  "123456",
		}
		body, _ := json.Marshal(req)

		// Return expired verification
		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnRows(sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
				AddRow("test@example.com", codeHash, time.Now().Add(-1*time.Hour), int(mail.MailConfirmation)))

		// Expect cleanup of expired verification
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostVerifyEmail() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "verification code has expired") {
			t.Errorf("PostVerifyEmail() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid verification code", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		// Hash for a different code
		codeHash, _ := auth.GetHash("654321")

		req := VerifyEmailRequest{
			Email: "test@example.com",
			Code:  "123456", // Wrong code
		}
		body, _ := json.Marshal(req)

		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnRows(sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
				AddRow("test@example.com", codeHash, time.Now().Add(1*time.Hour), int(mail.MailConfirmation)))

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostVerifyEmail() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "invalid verification code") {
			t.Errorf("PostVerifyEmail() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("verification not found", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := VerifyEmailRequest{
			Email: "test@example.com",
			Code:  "123456",
		}
		body, _ := json.Marshal(req)

		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnError(sql.ErrNoRows)

		httpReq := httptest.NewRequest("POST", "/user/verify-email", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostVerifyEmail(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostVerifyEmail() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "invalid or expired verification code") {
			t.Errorf("PostVerifyEmail() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestPostRequestPasswordReset(t *testing.T) {
	t.Run("successful password reset request - user exists", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockMS := &MockMailService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, MS: mockMS}

		req := PasswordResetRequest{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(req)

		// Expect GetUserByEmail query
		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
				AddRow(1, "testuser", "test@example.com", "hash123", true, 0))

		// Expect InsertVerification query
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs("test@example.com", int(mail.PassReset)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectExec("INSERT INTO verifications").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		httpReq := httptest.NewRequest("POST", "/user/password-reset", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostRequestPasswordReset(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("PostRequestPasswordReset() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		if !strings.Contains(response["message"], "password reset link has been sent") {
			t.Errorf("PostRequestPasswordReset() message = %s", response["message"])
		}

		// Verify email was sent
		if len(mockMS.SentMails) != 1 {
			t.Errorf("Expected 1 email to be sent, got %d", len(mockMS.SentMails))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("user not found - security response", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := PasswordResetRequest{
			Email: "nonexistent@example.com",
		}
		body, _ := json.Marshal(req)

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs("nonexistent@example.com").
			WillReturnError(sql.ErrNoRows)

		httpReq := httptest.NewRequest("POST", "/user/password-reset", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostRequestPasswordReset(w, httpReq)

		// Should return 200 for security reasons
		if w.Code != http.StatusOK {
			t.Errorf("PostRequestPasswordReset() status = %d, want %d", w.Code, http.StatusOK)
		}

		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		if !strings.Contains(response["message"], "password reset link has been sent") {
			t.Errorf("PostRequestPasswordReset() should not reveal user existence, got: %s", response["message"])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		httpReq := httptest.NewRequest("POST", "/user/password-reset", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostRequestPasswordReset(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostRequestPasswordReset() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing email field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := PasswordResetRequest{}
		body, _ := json.Marshal(req)

		httpReq := httptest.NewRequest("POST", "/user/password-reset", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostRequestPasswordReset(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostRequestPasswordReset() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "email is required") {
			t.Errorf("PostRequestPasswordReset() body = %s", w.Body.String())
		}
	})

	t.Run("email sending failure", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockMS := &MockMailService{ShouldFail: true}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, MS: mockMS}

		req := PasswordResetRequest{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(req)

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
				AddRow(1, "testuser", "test@example.com", "hash123", true, 0))

		httpReq := httptest.NewRequest("POST", "/user/password-reset", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostRequestPasswordReset(w, httpReq)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("PostRequestPasswordReset() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestPostResetPassword(t *testing.T) {
	t.Run("successful password reset", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		codeHash, _ := auth.GetHash("123456")

		req := ResetPasswordRequest{
			Email:       "test@example.com",
			Code:        "123456",
			NewPassword: "newpassword123",
		}
		body, _ := json.Marshal(req)

		// Expect GetVerification query
		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.PassReset)).
			WillReturnRows(sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
				AddRow("test@example.com", codeHash, time.Now().Add(1*time.Hour), int(mail.PassReset)))

		// Expect UpdateUserPassword query
		mock.ExpectExec("UPDATE users SET hash").
			WithArgs(sqlmock.AnyArg(), "test@example.com").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect DeleteVerification query
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs("test@example.com", int(mail.PassReset)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		httpReq := httptest.NewRequest("POST", "/user/reset-password", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResetPassword(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("PostResetPassword() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		if response["message"] != "password reset successfully" {
			t.Errorf("PostResetPassword() message = %s", response["message"])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		httpReq := httptest.NewRequest("POST", "/user/reset-password", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResetPassword(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostResetPassword() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := ResetPasswordRequest{
			Email: "test@example.com",
			Code:  "123456",
			// Missing NewPassword
		}
		body, _ := json.Marshal(req)

		httpReq := httptest.NewRequest("POST", "/user/reset-password", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResetPassword(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostResetPassword() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "email, code, and new_password are required") {
			t.Errorf("PostResetPassword() body = %s", w.Body.String())
		}
	})

	t.Run("expired verification code", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		codeHash, _ := auth.GetHash("123456")

		req := ResetPasswordRequest{
			Email:       "test@example.com",
			Code:        "123456",
			NewPassword: "newpassword123",
		}
		body, _ := json.Marshal(req)

		// Return expired verification
		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.PassReset)).
			WillReturnRows(sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
				AddRow("test@example.com", codeHash, time.Now().Add(-1*time.Hour), int(mail.PassReset)))

		// Expect cleanup
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs("test@example.com", int(mail.PassReset)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		httpReq := httptest.NewRequest("POST", "/user/reset-password", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResetPassword(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostResetPassword() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "verification code has expired") {
			t.Errorf("PostResetPassword() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid verification code", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		codeHash, _ := auth.GetHash("654321")

		req := ResetPasswordRequest{
			Email:       "test@example.com",
			Code:        "123456", // Wrong code
			NewPassword: "newpassword123",
		}
		body, _ := json.Marshal(req)

		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs("test@example.com", int(mail.PassReset)).
			WillReturnRows(sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
				AddRow("test@example.com", codeHash, time.Now().Add(1*time.Hour), int(mail.PassReset)))

		httpReq := httptest.NewRequest("POST", "/user/reset-password", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResetPassword(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostResetPassword() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "invalid verification code") {
			t.Errorf("PostResetPassword() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestPostResendVerification(t *testing.T) {
	t.Run("successful verification resend - unverified user", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockMS := &MockMailService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, MS: mockMS}

		req := PasswordResetRequest{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(req)

		// Expect GetUserByEmail query
		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
				AddRow(1, "testuser", "test@example.com", "hash123", false, 0))

		// Expect InsertVerification query
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs("test@example.com", int(mail.MailConfirmation)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectExec("INSERT INTO verifications").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		httpReq := httptest.NewRequest("POST", "/user/resend-verification", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResendVerification(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("PostResendVerification() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		if response["message"] != "verification email sent" {
			t.Errorf("PostResendVerification() message = %s", response["message"])
		}

		// Verify email was sent
		if len(mockMS.SentMails) != 1 {
			t.Errorf("Expected 1 email to be sent, got %d", len(mockMS.SentMails))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("user already verified", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := PasswordResetRequest{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(req)

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
				AddRow(1, "testuser", "test@example.com", "hash123", true, 0))

		httpReq := httptest.NewRequest("POST", "/user/resend-verification", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResendVerification(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostResendVerification() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "email is already verified") {
			t.Errorf("PostResendVerification() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := PasswordResetRequest{
			Email: "nonexistent@example.com",
		}
		body, _ := json.Marshal(req)

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs("nonexistent@example.com").
			WillReturnError(sql.ErrNoRows)

		httpReq := httptest.NewRequest("POST", "/user/resend-verification", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResendVerification(w, httpReq)

		if w.Code != http.StatusNotFound {
			t.Errorf("PostResendVerification() status = %d, want %d", w.Code, http.StatusNotFound)
		}
		if !strings.Contains(w.Body.String(), "user not found") {
			t.Errorf("PostResendVerification() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		httpReq := httptest.NewRequest("POST", "/user/resend-verification", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResendVerification(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostResendVerification() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing email field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := PasswordResetRequest{}
		body, _ := json.Marshal(req)

		httpReq := httptest.NewRequest("POST", "/user/resend-verification", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostResendVerification(w, httpReq)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostResendVerification() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "email is required") {
			t.Errorf("PostResendVerification() body = %s", w.Body.String())
		}
	})
}
