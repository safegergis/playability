package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"playability/db"
	"playability/types"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestGetFeatureReportsHandler(t *testing.T) {
	t.Run("successful feature reports calculation", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind", "full_controller_support", "controller_remapping"}).
			AddRow("true", "true", "true", "false").
			AddRow("true", "limited", "true", "true").
			AddRow("limited", "true", "false", "true")

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports").
			WithArgs(123).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/features/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetFeatureReportsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetFeatureReportsHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify response is valid JSON array
		var features []types.FeatureStat
		err = json.NewDecoder(w.Body).Decode(&features)
		if err != nil {
			t.Errorf("GetFeatureReportsHandler() returned invalid JSON: %v", err)
		}

		if len(features) != 4 {
			t.Errorf("GetFeatureReportsHandler() returned %d features, want 4", len(features))
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

		req := httptest.NewRequest("GET", "/features/invalid", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetFeatureReportsHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("GetFeatureReportsHandler() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("database error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports").
			WithArgs(123).
			WillReturnError(errors.New("database error"))

		req := httptest.NewRequest("GET", "/features/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetFeatureReportsHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetFeatureReportsHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
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

		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind", "full_controller_support", "controller_remapping"})

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports").
			WithArgs(999).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/features/999", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetFeatureReportsHandler(w, req)

		if w.Code != http.StatusNotAcceptable {
			t.Errorf("GetFeatureReportsHandler() status = %d, want %d", w.Code, http.StatusNotAcceptable)
		}

		if !strings.Contains(w.Body.String(), "Not enough reports") {
			t.Errorf("GetFeatureReportsHandler() body = %s, want 'Not enough reports'", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("single report returns valid features", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind", "full_controller_support", "controller_remapping"}).
			AddRow("true", "limited", "false", "true")

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports").
			WithArgs(456).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/features/456", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "456")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetFeatureReportsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetFeatureReportsHandler() status = %d, want %d", w.Code, http.StatusOK)
		}

		var features []types.FeatureStat
		json.NewDecoder(w.Body).Decode(&features)

		if len(features) != 4 {
			t.Errorf("GetFeatureReportsHandler() returned %d features, want 4", len(features))
		}

		// Check consensus values (single report, so percentages should be 1.0)
		for _, feature := range features {
			total := feature.TruePercentage + feature.LimitedPercentage + feature.FalsePercentage
			if total != 1.0 {
				t.Errorf("Feature %s: percentages sum to %f, want 1.0", feature.FeatureName, total)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}
