package db

import (
	"database/sql"
	"errors"
	"playability/pkg/mail"
	"playability/types"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertVerification(t *testing.T) {
	t.Run("successful insertion", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		verification := types.VerificationRow{
			Email:     "test@example.com",
			Code:      "123456",
			ExpiresAt: time.Now().Add(15 * time.Minute),
			Type:      mail.MailConfirmation,
		}

		// Expect delete old verifications
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs(verification.Email, verification.Type).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Expect insert new verification
		mock.ExpectExec("INSERT INTO verifications").
			WithArgs(verification.Email, verification.Code, verification.ExpiresAt, verification.Type).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = dbModel.InsertVerification(verification)
		if err != nil {
			t.Errorf("InsertVerification() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		verification := types.VerificationRow{
			Email:     "test@example.com",
			Code:      "123456",
			ExpiresAt: time.Now().Add(15 * time.Minute),
			Type:      mail.MailConfirmation,
		}

		err := dbModel.InsertVerification(verification)
		if err == nil {
			t.Error("InsertVerification() expected error for nil DB")
		}
	})

	t.Run("delete error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		verification := types.VerificationRow{
			Email:     "test@example.com",
			Code:      "123456",
			ExpiresAt: time.Now().Add(15 * time.Minute),
			Type:      mail.MailConfirmation,
		}

		// Expect delete to fail
		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs(verification.Email, verification.Type).
			WillReturnError(errors.New("delete error"))

		err = dbModel.InsertVerification(verification)
		if err == nil {
			t.Error("InsertVerification() expected error for delete failure")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestGetVerification(t *testing.T) {
	t.Run("verification found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"
		mailType := int(mail.MailConfirmation)
		expiresAt := time.Now().Add(15 * time.Minute)

		rows := sqlmock.NewRows([]string{"email", "code", "expires_at", "type"}).
			AddRow(email, "123456", expiresAt, mailType)

		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs(email, mailType).
			WillReturnRows(rows)

		verification, err := dbModel.GetVerification(email, mailType)
		if err != nil {
			t.Errorf("GetVerification() error = %v", err)
		}
		if verification.Email != email {
			t.Errorf("GetVerification() Email = %s, want %s", verification.Email, email)
		}
		if verification.Code != "123456" {
			t.Errorf("GetVerification() Code = %s, want '123456'", verification.Code)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("verification not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"
		mailType := int(mail.MailConfirmation)

		mock.ExpectQuery("SELECT email, code, expires_at, type FROM verifications").
			WithArgs(email, mailType).
			WillReturnError(sql.ErrNoRows)

		_, err = dbModel.GetVerification(email, mailType)
		if err == nil {
			t.Error("GetVerification() expected error for not found")
		}
		if err.Error() != "verification not found" {
			t.Errorf("GetVerification() error = %v, want 'verification not found'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		_, err := dbModel.GetVerification("test@example.com", int(mail.MailConfirmation))

		if err == nil {
			t.Error("GetVerification() expected error for nil DB")
		}
	})
}

func TestDeleteVerification(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"
		mailType := int(mail.MailConfirmation)

		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs(email, mailType).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = dbModel.DeleteVerification(email, mailType)
		if err != nil {
			t.Errorf("DeleteVerification() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("verification not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"
		mailType := int(mail.MailConfirmation)

		mock.ExpectExec("DELETE FROM verifications WHERE email").
			WithArgs(email, mailType).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = dbModel.DeleteVerification(email, mailType)
		if err == nil {
			t.Error("DeleteVerification() expected error for not found")
		}
		if err.Error() != "verification not found" {
			t.Errorf("DeleteVerification() error = %v, want 'verification not found'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		err := dbModel.DeleteVerification("test@example.com", int(mail.MailConfirmation))

		if err == nil {
			t.Error("DeleteVerification() expected error for nil DB")
		}
	})
}

func TestCleanExpiredVerifications(t *testing.T) {
	t.Run("successful cleanup", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}

		mock.ExpectExec("DELETE FROM verifications WHERE expires_at").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err = dbModel.CleanExpiredVerifications()
		if err != nil {
			t.Errorf("CleanExpiredVerifications() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		err := dbModel.CleanExpiredVerifications()

		if err == nil {
			t.Error("CleanExpiredVerifications() expected error for nil DB")
		}
	})
}

func TestMarkUserAsVerified(t *testing.T) {
	t.Run("successful verification", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"

		mock.ExpectExec("UPDATE users SET verified").
			WithArgs(email).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = dbModel.MarkUserAsVerified(email)
		if err != nil {
			t.Errorf("MarkUserAsVerified() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "nonexistent@example.com"

		mock.ExpectExec("UPDATE users SET verified").
			WithArgs(email).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = dbModel.MarkUserAsVerified(email)
		if err == nil {
			t.Error("MarkUserAsVerified() expected error for not found")
		}
		if err.Error() != "user not found" {
			t.Errorf("MarkUserAsVerified() error = %v, want 'user not found'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		err := dbModel.MarkUserAsVerified("test@example.com")

		if err == nil {
			t.Error("MarkUserAsVerified() expected error for nil DB")
		}
	})
}

func TestUpdateUserPassword(t *testing.T) {
	t.Run("successful password update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"
		newHash := "newhash123"

		mock.ExpectExec("UPDATE users SET hash").
			WithArgs(newHash, email).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = dbModel.UpdateUserPassword(email, newHash)
		if err != nil {
			t.Errorf("UpdateUserPassword() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "nonexistent@example.com"
		newHash := "newhash123"

		mock.ExpectExec("UPDATE users SET hash").
			WithArgs(newHash, email).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = dbModel.UpdateUserPassword(email, newHash)
		if err == nil {
			t.Error("UpdateUserPassword() expected error for not found")
		}
		if err.Error() != "user not found" {
			t.Errorf("UpdateUserPassword() error = %v, want 'user not found'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		err := dbModel.UpdateUserPassword("test@example.com", "newhash")

		if err == nil {
			t.Error("UpdateUserPassword() expected error for nil DB")
		}
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("user found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"

		rows := sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
			AddRow(42, "testuser", email, "hash123", true, 5)

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs(email).
			WillReturnRows(rows)

		user, err := dbModel.GetUserByEmail(email)
		if err != nil {
			t.Errorf("GetUserByEmail() error = %v", err)
		}
		if user.ID != 42 {
			t.Errorf("GetUserByEmail() ID = %d, want 42", user.ID)
		}
		if user.Email != email {
			t.Errorf("GetUserByEmail() Email = %s, want %s", user.Email, email)
		}
		if !user.Verified {
			t.Error("GetUserByEmail() Verified = false, want true")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "nonexistent@example.com"

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE email").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		_, err = dbModel.GetUserByEmail(email)
		if err == nil {
			t.Error("GetUserByEmail() expected error for not found")
		}
		if err.Error() != "user not found" {
			t.Errorf("GetUserByEmail() error = %v, want 'user not found'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		_, err := dbModel.GetUserByEmail("test@example.com")

		if err == nil {
			t.Error("GetUserByEmail() expected error for nil DB")
		}
	})
}
