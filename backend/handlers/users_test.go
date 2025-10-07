package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"playability/auth"
	"playability/db"
	"playability/types"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestPostCreateUser(t *testing.T) {
	t.Run("successful user registration", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Username: "TestUser",
			Email:    "Test@Example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		// Expect check for existing user
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("test@example.com", "testuser").
			WillReturnError(sql.ErrNoRows)

		// Expect insert
		mock.ExpectExec("INSERT INTO users").
			WithArgs("testuser", "test@example.com", sqlmock.AnyArg(), 0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusCreated)
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

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "Invalid request body") {
			t.Errorf("PostCreateUser() body = %s, want 'Invalid request body'", w.Body.String())
		}
	})

	t.Run("missing username field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "username, email, and password are required") {
			t.Errorf("PostCreateUser() body = %s", w.Body.String())
		}
	})

	t.Run("missing email field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Username: "testuser",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing password field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Username: "testuser",
			Email:    "test@example.com",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("duplicate email error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Username: "testuser",
			Email:    "existing@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		// Expect check returns existing user
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("existing@example.com", "testuser").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Expect check for email specifically
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("existing@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusConflict)
		}
		if !strings.Contains(w.Body.String(), "email is already in use") {
			t.Errorf("PostCreateUser() body = %s, want 'email is already in use'", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("duplicate username error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Username: "existinguser",
			Email:    "new@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		// Expect check returns existing user
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("new@example.com", "existinguser").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Expect check for email (not found)
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("new@example.com").
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusConflict)
		}
		if !strings.Contains(w.Body.String(), "username is already in use") {
			t.Errorf("PostCreateUser() body = %s, want 'username is already in use'", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("database internal error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		// Expect check for existing user
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("test@example.com", "testuser").
			WillReturnError(errors.New("database error"))

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("lowercase normalization", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserRegister{
			Username: "TestUser",
			Email:    "Test@Example.COM",
			Password: "Password123",
		}
		body, _ := json.Marshal(user)

		// Expect lowercase values
		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("test@example.com", "testuser").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO users").
			WithArgs("testuser", "test@example.com", sqlmock.AnyArg(), 0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostCreateUser(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("PostCreateUser() status = %d, want %d", w.Code, http.StatusCreated)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestPostLoginUser(t *testing.T) {
	// Set JWT_SECRET for CreateToken
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Note: QueryUser doesn't select 'verified' column, so Verified is always false (zero value)
	// This means the check at line 137 never blocks anyone (double bug!)
	t.Run("successful login - verified check never triggers due to missing column", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserLogin{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		// Generate real hash for password
		hash, _ := auth.GetHash("password123")

		// Expect CheckUser queries
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"hash"}).AddRow(hash))

		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

		// Expect QueryUser
		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(42).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
				AddRow(42, "testuser", "test@example.com", hash, false, 5))

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("PostLoginUser() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify response contains token and user info
		var response types.LoginResponse
		json.NewDecoder(w.Body).Decode(&response)
		if response.Token == "" {
			t.Error("PostLoginUser() response missing token")
		}
		if response.User.ID != 42 {
			t.Errorf("PostLoginUser() user ID = %d, want 42", response.User.ID)
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

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostLoginUser() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing email field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserLogin{
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostLoginUser() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "email and password are required") {
			t.Errorf("PostLoginUser() body = %s", w.Body.String())
		}
	})

	t.Run("missing password field", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserLogin{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostLoginUser() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid credentials - user not found", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserLogin{
			Email:    "notfound@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		// Expect CheckUser query returns no rows
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("PostLoginUser() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
		if !strings.Contains(w.Body.String(), "invalid email or password") {
			t.Errorf("PostLoginUser() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid credentials - wrong password", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserLogin{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(user)

		// Generate hash for correct password
		hash, _ := auth.GetHash("correctpassword")

		// Expect CheckUser query
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"hash"}).AddRow(hash))

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("PostLoginUser() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	// Note: CheckUser returns (0, false, nil) on DB errors, treating them as invalid credentials
	t.Run("database error during CheckUser treated as invalid credentials", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserLogin{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		// Expect CheckUser query with error - returns (0, false, nil)
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnError(errors.New("database error"))

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		// CheckUser returns valid=false on DB error, so handler returns 401 not 500
		if w.Code != http.StatusUnauthorized {
			t.Errorf("PostLoginUser() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("database error during QueryUser", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		user := types.UserLogin{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(user)

		hash, _ := auth.GetHash("password123")

		// Expect CheckUser queries succeed
		mock.ExpectQuery("SELECT hash FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"hash"}).AddRow(hash))

		mock.ExpectQuery("SELECT id FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

		// Expect QueryUser fails
		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(42).
			WillReturnError(errors.New("database error"))

		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostLoginUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("PostLoginUser() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

}

func TestGetUserHandler(t *testing.T) {
	t.Run("successful user retrieval", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(42).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "hash", "verified", "num_reports"}).
				AddRow(42, "testuser", "test@example.com", "hash123", true, 5))

		req := httptest.NewRequest("GET", "/user/42", nil)
		w := httptest.NewRecorder()

		// Set up chi context with URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "42")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetUserHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetUserHandler() status = %d, want %d", w.Code, http.StatusOK)
		}

		var user types.UserRow
		json.NewDecoder(w.Body).Decode(&user)
		if user.ID != 42 {
			t.Errorf("GetUserHandler() user ID = %d, want 42", user.ID)
		}
		if user.Username != "testuser" {
			t.Errorf("GetUserHandler() username = %s, want 'testuser'", user.Username)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid user ID format", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := httptest.NewRequest("GET", "/user/invalid", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetUserHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("GetUserHandler() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if !strings.Contains(w.Body.String(), "Invalid user ID") {
			t.Errorf("GetUserHandler() body = %s", w.Body.String())
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest("GET", "/user/999", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetUserHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("GetUserHandler() status = %d, want %d", w.Code, http.StatusNotFound)
		}
		if !strings.Contains(w.Body.String(), "User not found") {
			t.Errorf("GetUserHandler() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		mock.ExpectQuery("SELECT id, username, email, hash, verified, num_reports FROM users WHERE id").
			WithArgs(42).
			WillReturnError(errors.New("database error"))

		req := httptest.NewRequest("GET", "/user/42", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "42")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetUserHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("GetUserHandler() status = %d, want %d", w.Code, http.StatusNotFound)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}
