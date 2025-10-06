package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"playability/db"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestGetScoreHandler(t *testing.T) {
	t.Run("successful score calculation", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		rows := sqlmock.NewRows([]string{"score"}).
			AddRow(8).
			AddRow(7).
			AddRow(9).
			AddRow(6)

		mock.ExpectQuery("SELECT score FROM reports WHERE game_id").
			WithArgs(123).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/score/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetScoreHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetScoreHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify response is a number
		body := strings.TrimSpace(w.Body.String())
		if body == "" {
			t.Error("GetScoreHandler() returned empty body")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("invalid game ID format", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		req := httptest.NewRequest("GET", "/score/invalid", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetScoreHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetScoreHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("database error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		mock.ExpectQuery("SELECT score FROM reports WHERE game_id").
			WithArgs(123).
			WillReturnError(errors.New("database error"))

		req := httptest.NewRequest("GET", "/score/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetScoreHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetScoreHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("not enough reports - empty results", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		rows := sqlmock.NewRows([]string{"score"})

		mock.ExpectQuery("SELECT score FROM reports WHERE game_id").
			WithArgs(999).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/score/999", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetScoreHandler(w, req)

		if w.Code != http.StatusNotAcceptable {
			t.Errorf("GetScoreHandler() status = %d, want %d", w.Code, http.StatusNotAcceptable)
		}

		if !strings.Contains(w.Body.String(), "Not enough reports") {
			t.Errorf("GetScoreHandler() body = %s, want 'Not enough reports'", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("single score returns valid average", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		rows := sqlmock.NewRows([]string{"score"}).AddRow(8)

		mock.ExpectQuery("SELECT score FROM reports WHERE game_id").
			WithArgs(456).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/score/456", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "456")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetScoreHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetScoreHandler() status = %d, want %d", w.Code, http.StatusOK)
		}

		body := strings.TrimSpace(w.Body.String())
		if body != "8.000000" {
			t.Errorf("GetScoreHandler() body = %s, want '8.000000'", body)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}
