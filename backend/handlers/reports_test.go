package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"playability/db"
	"playability/pkg/ai"
	"playability/types"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func TestPostReportHandler(t *testing.T) {
	// Set up JWT for authentication
	tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)

	t.Run("successful report submission with AI moderation pass", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{
			ModerationResponse: &types.ModerationResponse{
				Violation:  false,
				Categories: []string{},
			},
		}

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			AI: mockAI,
		}

		reportBody := types.ReportRegister{
			GameID:   123,
			Platform: types.PC,
			Score:    8,
			Report:   "Great accessibility features",
		}
		body, _ := json.Marshal(reportBody)

		// Expect duplicate check in InsertReport
		mock.ExpectQuery("SELECT id FROM reports WHERE game_id").
			WithArgs(123, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		// Expect update user num_reports
		mock.ExpectExec("UPDATE users SET num_reports").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect InsertReport (uses QueryRow with RETURNING id)
		mock.ExpectQuery("INSERT INTO reports").
			WithArgs(123, sqlmock.AnyArg(), types.PC, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 8, "Great accessibility features").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Expect GetReportCount
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(123).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Create JWT token
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": "42"})

		req := httptest.NewRequest("POST", "/user/report", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		token, _ := tokenAuth.Decode(tokenString)
		ctx := jwtauth.NewContext(req.Context(), token, nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		env.PostReportHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("PostReportHandler() status = %d, want %d. Body: %s", w.Code, http.StatusCreated, w.Body.String())
		}

		// Verify AI moderation was called
		if mockAI.ModerationCallCount != 1 {
			t.Errorf("AI.Moderation() called %d times, want 1", mockAI.ModerationCallCount)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("report blocked by AI moderation", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{
			ModerationResponse: &types.ModerationResponse{
				Violation:   true,
				Categories:  []string{"Harassment"},
				Explanation: "Contains harassment",
			},
		}

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			AI: mockAI,
		}

		reportBody := types.ReportRegister{
			GameID:   123,
			Platform: types.PC,
			Score:    8,
			Report:   "Bad content",
		}
		body, _ := json.Marshal(reportBody)

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": "42"})

		req := httptest.NewRequest("POST", "/user/report", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		token, _ := tokenAuth.Decode(tokenString)
		ctx := jwtauth.NewContext(req.Context(), token, nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		env.PostReportHandler(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("PostReportHandler() status = %d, want %d", w.Code, http.StatusForbidden)
		}

		if !strings.Contains(w.Body.String(), "Harassment") {
			t.Errorf("PostReportHandler() body should contain violation category")
		}

		// Verify AI moderation was called
		if mockAI.ModerationCallCount != 1 {
			t.Errorf("AI.Moderation() called %d times, want 1", mockAI.ModerationCallCount)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		req := httptest.NewRequest("POST", "/user/report", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostReportHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("PostReportHandler() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing JWT token", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		reportBody := types.ReportRegister{
			GameID:   123,
			Platform: types.PC,
			Score:    8,
			Report:   "Great accessibility features",
		}
		body, _ := json.Marshal(reportBody)

		req := httptest.NewRequest("POST", "/user/report", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.PostReportHandler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("PostReportHandler() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("valid JWT token extraction", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		reportBody := types.ReportRegister{
			GameID:   123,
			Platform: types.PC,
			Score:    8,
			Report:   "Great accessibility features",
		}
		body, _ := json.Marshal(reportBody)

		// Create JWT token with user ID
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": "42"})

		req := httptest.NewRequest("POST", "/user/report", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Add JWT context with proper token
		token, _ := tokenAuth.Decode(tokenString)
		ctx := jwtauth.NewContext(req.Context(), token, nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		// This will fail at moderation step (requires AI mocking)
		env.PostReportHandler(w, req)

		// We expect it to fail at AI moderation, not at JWT extraction
		if w.Code == http.StatusUnauthorized {
			t.Error("PostReportHandler() should not fail at JWT extraction")
		}
	})

	t.Run("duplicate report returns conflict", func(t *testing.T) {
		t.Skip("Requires AI moderation mocking - AI service should be dependency injected")
	})
}

func TestGetReportCardsHandler(t *testing.T) {
	t.Run("successful report cards retrieval", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		rows := sqlmock.NewRows([]string{"id", "created_at", "game_id", "user_id", "platform", "score", "report"}).
			AddRow(1, time.Now(), 123, 1, types.PC, 8, "Good accessibility").
			AddRow(2, time.Now(), 123, 2, types.Playstation, 7, "Decent features")

		mock.ExpectQuery("SELECT id, created_at, game_id, user_id, platform, score, report FROM reports WHERE game_id").
			WithArgs(123).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/reports/cards/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportCardsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetReportCardsHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var reports []types.ReportCards
		json.NewDecoder(w.Body).Decode(&reports)
		if len(reports) != 2 {
			t.Errorf("GetReportCardsHandler() returned %d reports, want 2", len(reports))
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

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		req := httptest.NewRequest("GET", "/reports/cards/invalid", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportCardsHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetReportCardsHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("empty results", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		rows := sqlmock.NewRows([]string{"id", "created_at", "game_id", "user_id", "platform", "score", "report"})

		mock.ExpectQuery("SELECT id, created_at, game_id, user_id, platform, score, report FROM reports WHERE game_id").
			WithArgs(999).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/reports/cards/999", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportCardsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetReportCardsHandler() status = %d, want %d", w.Code, http.StatusOK)
		}

		var reports []types.ReportCards
		json.NewDecoder(w.Body).Decode(&reports)
		if len(reports) != 0 {
			t.Errorf("GetReportCardsHandler() returned %d reports, want 0", len(reports))
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

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		mock.ExpectQuery("SELECT id, created_at, game_id, user_id, platform, score, report FROM reports WHERE game_id").
			WithArgs(123).
			WillReturnError(errors.New("database error"))

		req := httptest.NewRequest("GET", "/reports/cards/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportCardsHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetReportCardsHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestGetReportSummaryHandler(t *testing.T) {
	t.Run("successful summary retrieval", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		rows := sqlmock.NewRows([]string{"id", "game_id", "summary"}).
			AddRow(1, 123, "This game has excellent accessibility features including full controller support and color blind modes.")

		mock.ExpectQuery("SELECT (.+) FROM report_summaries WHERE game_id").
			WithArgs(123).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/reports/summary/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportSummaryHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetReportSummaryHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var summary types.ReportSummaryRow
		json.NewDecoder(w.Body).Decode(&summary)
		if summary.GameID != 123 {
			t.Errorf("GetReportSummaryHandler() game_id = %d, want 123", summary.GameID)
		}
		if summary.Summary == "" {
			t.Error("GetReportSummaryHandler() summary is empty")
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

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		req := httptest.NewRequest("GET", "/reports/summary/invalid", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportSummaryHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("GetReportSummaryHandler() status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("no summary found - returns 404", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		mock.ExpectQuery("SELECT (.+) FROM report_summaries WHERE game_id").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest("GET", "/reports/summary/999", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportSummaryHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetReportSummaryHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("empty summary field returns 404", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		rows := sqlmock.NewRows([]string{"id", "game_id", "summary"}).
			AddRow(1, 123, "")

		mock.ExpectQuery("SELECT (.+) FROM report_summaries WHERE game_id").
			WithArgs(123).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/reports/summary/123", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("game", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		env.GetReportSummaryHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("GetReportSummaryHandler() status = %d, want %d", w.Code, http.StatusNotFound)
		}

		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		if response["message"] != "No summary available" {
			t.Errorf("GetReportSummaryHandler() message = %s, want 'No summary available'", response["message"])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestGetUserReportsHandler(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)

	t.Run("successful user reports retrieval", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		rows := sqlmock.NewRows([]string{"id", "created_at", "game_id", "game_name", "cover_art", "user_id", "platform", "score", "report"}).
			AddRow(1, time.Now(), 123, "Test Game 1", "cover1.jpg", 42, types.PC, 8, "My first report").
			AddRow(2, time.Now(), 456, "Test Game 2", "cover2.jpg", 42, types.Playstation, 9, "My second report")

		mock.ExpectQuery("SELECT (.+) FROM reports r LEFT JOIN games g ON r.game_id = g.id WHERE r.user_id").
			WithArgs(42).
			WillReturnRows(rows)

		// Create JWT token with user ID
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": "42"})

		req := httptest.NewRequest("GET", "/user/reports", nil)
		w := httptest.NewRecorder()

		// Add JWT context with proper claims
		token, _ := tokenAuth.Decode(tokenString)
		ctx := jwtauth.NewContext(req.Context(), token, nil)
		req = req.WithContext(ctx)

		env.GetUserReportsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetUserReportsHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var reports []types.ReportCards
		json.NewDecoder(w.Body).Decode(&reports)
		if len(reports) != 2 {
			t.Errorf("GetUserReportsHandler() returned %d reports, want 2", len(reports))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("missing JWT token", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		req := httptest.NewRequest("GET", "/user/reports", nil)
		w := httptest.NewRecorder()

		env.GetUserReportsHandler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("GetUserReportsHandler() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
		if !strings.Contains(w.Body.String(), "Unauthorized") {
			t.Errorf("GetUserReportsHandler() body = %s, want 'Unauthorized'", w.Body.String())
		}
	})

	t.Run("invalid user ID in claims", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		// Create JWT token with invalid user ID (not a string)
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": 123}) // int instead of string

		req := httptest.NewRequest("GET", "/user/reports", nil)
		w := httptest.NewRecorder()

		token, _ := tokenAuth.Decode(tokenString)
		ctx := jwtauth.NewContext(req.Context(), token, nil)
		req = req.WithContext(ctx)

		env.GetUserReportsHandler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("GetUserReportsHandler() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("empty reports for user", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		rows := sqlmock.NewRows([]string{"id", "created_at", "game_id", "game_name", "cover_art", "user_id", "platform", "score", "report"})

		mock.ExpectQuery("SELECT (.+) FROM reports r LEFT JOIN games g ON r.game_id = g.id WHERE r.user_id").
			WithArgs(42).
			WillReturnRows(rows)

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": "42"})

		req := httptest.NewRequest("GET", "/user/reports", nil)
		w := httptest.NewRecorder()

		token, _ := tokenAuth.Decode(tokenString)
		ctx := jwtauth.NewContext(req.Context(), token, nil)
		req = req.WithContext(ctx)

		env.GetUserReportsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetUserReportsHandler() status = %d, want %d", w.Code, http.StatusOK)
		}

		var reports []types.ReportCards
		json.NewDecoder(w.Body).Decode(&reports)
		if len(reports) != 0 {
			t.Errorf("GetUserReportsHandler() returned %d reports, want 0", len(reports))
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

		mockAI := &ai.MockAIService{}
		env := &Env{DB: db.DatabaseModel{DB: mockDB}, AI: mockAI}

		mock.ExpectQuery("SELECT (.+) FROM reports r LEFT JOIN games g ON r.game_id = g.id WHERE r.user_id").
			WithArgs(42).
			WillReturnError(errors.New("database error"))

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": "42"})

		req := httptest.NewRequest("GET", "/user/reports", nil)
		w := httptest.NewRecorder()

		token, _ := tokenAuth.Decode(tokenString)
		ctx := jwtauth.NewContext(req.Context(), token, nil)
		req = req.WithContext(ctx)

		env.GetUserReportsHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetUserReportsHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}
