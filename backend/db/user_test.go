package db

import (
	"database/sql"
	"errors"
	"playability/auth"
	"playability/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestInsertUser(t *testing.T) {
	t.Run("successful user insertion", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		user := types.UserRegister{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Expect check for existing user
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs(user.Email, user.Username).
			WillReturnError(sql.ErrNoRows)

		// Expect insert
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.Username, user.Email, sqlmock.AnyArg(), 0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = dbModel.InsertUser(user)
		if err != nil {
			t.Errorf("InsertUser() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("duplicate email error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		user := types.UserRegister{
			Username: "testuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		// Expect check returns existing email
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs(user.Email, user.Username).
			WillReturnRows(rows)

		// Expect check for email specifically
		emailRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs(user.Email).
			WillReturnRows(emailRows)

		err = dbModel.InsertUser(user)
		if err == nil {
			t.Error("InsertUser() expected error for duplicate email")
		}
		if err.Error() != "email is already in use" {
			t.Errorf("InsertUser() error = %v, want 'email is already in use'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("duplicate username error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		user := types.UserRegister{
			Username: "existinguser",
			Email:    "new@example.com",
			Password: "password123",
		}

		// Expect check returns existing user
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs(user.Email, user.Username).
			WillReturnRows(rows)

		// Expect check for email (not found)
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs(user.Email).
			WillReturnError(sql.ErrNoRows)

		err = dbModel.InsertUser(user)
		if err == nil {
			t.Error("InsertUser() expected error for duplicate username")
		}
		if err.Error() != "username is already in use" {
			t.Errorf("InsertUser() error = %v, want 'username is already in use'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("unique constraint violation on insert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		user := types.UserRegister{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Expect check for existing user (none found)
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs(user.Email, user.Username).
			WillReturnError(sql.ErrNoRows)

		// Expect insert fails with unique violation
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.Username, user.Email, sqlmock.AnyArg(), 0).
			WillReturnError(&pq.Error{Code: "23505", Constraint: "unique_email"})

		err = dbModel.InsertUser(user)
		if err == nil {
			t.Error("InsertUser() expected error for unique constraint violation")
		}
		if err.Error() != "email is already in use" {
			t.Errorf("InsertUser() error = %v, want 'email is already in use'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		user := types.UserRegister{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		err := dbModel.InsertUser(user)
		if err == nil {
			t.Error("InsertUser() expected error for nil DB")
		}
		if err.Error() != "database connection is nil" {
			t.Errorf("InsertUser() error = %v, want 'database connection is nil'", err)
		}
	})
}

func TestCheckUser(t *testing.T) {
	t.Run("valid credentials", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"
		password := "password123"

		// Generate actual hash using auth package
		hash, err := auth.GetHash(password)
		if err != nil {
			t.Fatalf("Failed to generate hash: %v", err)
		}

		// Expect hash query
		hashRows := sqlmock.NewRows([]string{"hash"}).AddRow(hash)
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs(email).
			WillReturnRows(hashRows)

		// Expect ID query
		idRows := sqlmock.NewRows([]string{"id"}).AddRow(42)
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs(email).
			WillReturnRows(idRows)

		id, valid, err := dbModel.CheckUser(email, password)
		if err != nil {
			t.Errorf("CheckUser() error = %v", err)
		}
		if !valid {
			t.Error("CheckUser() valid = false, want true")
		}
		if id != 42 {
			t.Errorf("CheckUser() id = %d, want 42", id)
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
		password := "password123"

		// Expect hash query returns no rows
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		id, valid, err := dbModel.CheckUser(email, password)
		if err != nil {
			t.Errorf("CheckUser() error = %v, want nil", err)
		}
		if valid {
			t.Error("CheckUser() valid = true, want false")
		}
		if id != 0 {
			t.Errorf("CheckUser() id = %d, want 0", id)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		email := "test@example.com"
		correctPassword := "password123"
		wrongPassword := "wrongpassword"

		// Generate hash for correct password
		hash, err := auth.GetHash(correctPassword)
		if err != nil {
			t.Fatalf("Failed to generate hash: %v", err)
		}

		// Expect hash query
		hashRows := sqlmock.NewRows([]string{"hash"}).AddRow(hash)
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs(email).
			WillReturnRows(hashRows)

		id, valid, err := dbModel.CheckUser(email, wrongPassword)
		if err != nil {
			t.Errorf("CheckUser() error = %v, want nil", err)
		}
		if valid {
			t.Error("CheckUser() valid = true, want false")
		}
		if id != 0 {
			t.Errorf("CheckUser() id = %d, want 0", id)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		id, valid, err := dbModel.CheckUser("test@example.com", "password")

		if err == nil {
			t.Error("CheckUser() expected error for nil DB")
		}
		if valid {
			t.Error("CheckUser() valid = true, want false")
		}
		if id != 0 {
			t.Errorf("CheckUser() id = %d, want 0", id)
		}
	})
}

func TestQueryUser(t *testing.T) {
	t.Run("user found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		userID := 42

		rows := sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
			AddRow(42, "testuser", "test@example.com", "hash123", true, 5)

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(userID).
			WillReturnRows(rows)

		user, err := dbModel.QueryUser(userID)
		if err != nil {
			t.Errorf("QueryUser() error = %v", err)
		}
		if user.ID != 42 {
			t.Errorf("QueryUser() ID = %d, want 42", user.ID)
		}
		if user.Username != "testuser" {
			t.Errorf("QueryUser() Username = %s, want 'testuser'", user.Username)
		}
		if user.Email != "test@example.com" {
			t.Errorf("QueryUser() Email = %s, want 'test@example.com'", user.Email)
		}
		if user.NumOfReports != 5 {
			t.Errorf("QueryUser() NumOfReports = %d, want 5", user.NumOfReports)
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
		userID := 999

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := dbModel.QueryUser(userID)
		if err == nil {
			t.Error("QueryUser() expected error for non-existent user")
		}
		if user.ID != 0 {
			t.Errorf("QueryUser() returned non-empty user for error case")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := DatabaseModel{DB: db}
		userID := 42

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(userID).
			WillReturnError(errors.New("database error"))

		_, err = dbModel.QueryUser(userID)
		if err == nil {
			t.Error("QueryUser() expected error for database error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := DatabaseModel{DB: nil}
		_, err := dbModel.QueryUser(42)

		if err == nil {
			t.Error("QueryUser() expected error for nil DB")
		}
		if err.Error() != "database connection is nil" {
			t.Errorf("QueryUser() error = %v, want 'database connection is nil'", err)
		}
	})
}
