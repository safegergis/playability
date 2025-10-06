package db

import (
	"database/sql"
	"playability/types"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertReport(t *testing.T) {
	t.Run("successful report insertion", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		report := &types.ReportRow{
			GameID:                123,
			UserID:                42,
			Platform:              types.PC,
			ClosedCaptions:        "true",
			ColorBlind:            "limited",
			FullControllerSupport: "true",
			ControllerRemapping:   "false",
			Score:                 8,
			Report:                "Great accessibility features!",
		}

		// Expect check for existing report
		mock.ExpectQuery("SELECT id FROM reports WHERE game_id").
			WithArgs(report.GameID, report.UserID).
			WillReturnError(sql.ErrNoRows)

		// Expect update user num_reports
		mock.ExpectExec("UPDATE users SET num_reports").
			WithArgs(report.UserID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect insert
		mock.ExpectQuery("INSERT INTO reports").
			WithArgs(report.GameID, report.UserID, report.Platform, report.ClosedCaptions,
				report.ColorBlind, report.FullControllerSupport, report.ControllerRemapping,
				report.Score, report.Report).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err = dbModel.InsertReport(report)
		if err != nil {
			t.Errorf("InsertReport() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("duplicate report error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		report := &types.ReportRow{
			GameID: 123,
			UserID: 42,
		}

		// Expect check returns existing report
		rows := sqlmock.NewRows([]string{"id"}).AddRow(999)
		mock.ExpectQuery("SELECT id FROM reports WHERE game_id").
			WithArgs(report.GameID, report.UserID).
			WillReturnRows(rows)

		err = dbModel.InsertReport(report)
		if err == nil {
			t.Error("InsertReport() expected error for duplicate report")
		}
		if err.Error() != "report already exists" {
			t.Errorf("InsertReport() error = %v, want 'report already exists'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil report", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}

		err = dbModel.InsertReport(nil)
		if err == nil {
			t.Error("InsertReport() expected error for nil report")
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		report := &types.ReportRow{GameID: 123, UserID: 42}

		err := dbModel.InsertReport(report)
		if err == nil {
			t.Error("InsertReport() expected error for nil DB")
		}
	})
}

func TestInsertReportSummary(t *testing.T) {
	t.Run("successful summary insertion", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		summary := &types.ReportSummaryRow{
			GameID:  123,
			Summary: "Overall, the game has excellent accessibility features.",
		}

		mock.ExpectExec("INSERT INTO report_summaries").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = dbModel.InsertReportSummary(summary)
		if err != nil {
			t.Errorf("InsertReportSummary() error = %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil summary", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}

		err = dbModel.InsertReportSummary(nil)
		if err == nil {
			t.Error("InsertReportSummary() expected error for nil summary")
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		summary := &types.ReportSummaryRow{GameID: 123}

		err := dbModel.InsertReportSummary(summary)
		if err == nil {
			t.Error("InsertReportSummary() expected error for nil DB")
		}
	})
}

func TestQueryReportCards(t *testing.T) {
	t.Run("reports found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "created_at", "game_id", "user_id", "platform", "score", "report"}).
			AddRow(1, now, 123, 42, types.PC, 8, "Great game!").
			AddRow(2, now, 123, 43, types.Playstation, 7, "Good accessibility")

		mock.ExpectQuery("SELECT id, created_at, game_id, user_id, platform, score, report FROM reports").
			WithArgs(gameID).
			WillReturnRows(rows)

		reports, err := dbModel.QueryReportCards(gameID)
		if err != nil {
			t.Errorf("QueryReportCards() error = %v", err)
		}
		if len(reports) != 2 {
			t.Errorf("QueryReportCards() returned %d reports, want 2", len(reports))
		}
		if reports[0].Score != 8 {
			t.Errorf("QueryReportCards() first report score = %d, want 8", reports[0].Score)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		_, err := dbModel.QueryReportCards(123)

		if err == nil {
			t.Error("QueryReportCards() expected error for nil DB")
		}
	})
}

func TestQueryAccessibilityScores(t *testing.T) {
	t.Run("scores found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		rows := sqlmock.NewRows([]string{"score"}).
			AddRow(8).
			AddRow(7).
			AddRow(9)

		mock.ExpectQuery("SELECT score FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		scores, err := dbModel.QueryAccessibilityScores(gameID)
		if err != nil {
			t.Errorf("QueryAccessibilityScores() error = %v", err)
		}
		if len(scores) != 3 {
			t.Errorf("QueryAccessibilityScores() returned %d scores, want 3", len(scores))
		}
		if scores[0] != 8 || scores[1] != 7 || scores[2] != 9 {
			t.Errorf("QueryAccessibilityScores() scores = %v, want [8 7 9]", scores)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("no scores", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 999

		rows := sqlmock.NewRows([]string{"score"})
		mock.ExpectQuery("SELECT score FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		scores, err := dbModel.QueryAccessibilityScores(gameID)
		if err != nil {
			t.Errorf("QueryAccessibilityScores() error = %v", err)
		}
		if len(scores) != 0 {
			t.Errorf("QueryAccessibilityScores() returned %d scores, want 0", len(scores))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		_, err := dbModel.QueryAccessibilityScores(123)

		if err == nil {
			t.Error("QueryAccessibilityScores() expected error for nil DB")
		}
	})
}

func TestGetReportCount(t *testing.T) {
	t.Run("count found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		rows := sqlmock.NewRows([]string{"count"}).AddRow(10)
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(gameID).
			WillReturnRows(rows)

		count, err := dbModel.GetReportCount(gameID)
		if err != nil {
			t.Errorf("GetReportCount() error = %v", err)
		}
		if count != 10 {
			t.Errorf("GetReportCount() = %d, want 10", count)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		_, err := dbModel.GetReportCount(123)

		if err == nil {
			t.Error("GetReportCount() expected error for nil DB")
		}
	})
}

func TestQueryReportsForSummarization(t *testing.T) {
	t.Run("reports found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		rows := sqlmock.NewRows([]string{"id", "game_id", "user_id", "platform", "closed_captions", "color_blind", "full_controller_support", "controller_remapping", "score", "report"}).
			AddRow(1, 123, 42, types.PC, "true", "true", "true", "false", 8, "Great features").
			AddRow(2, 123, 43, types.Xbox, "true", "limited", "true", "true", 7, "Good support")

		mock.ExpectQuery("SELECT id, game_id, user_id, platform, closed_captions, color_blind, full_controller_support, controller_remapping, score, report FROM reports").
			WithArgs(gameID).
			WillReturnRows(rows)

		reports, err := dbModel.QueryReportsForSummarization(gameID)
		if err != nil {
			t.Errorf("QueryReportsForSummarization() error = %v", err)
		}
		if len(reports) != 2 {
			t.Errorf("QueryReportsForSummarization() returned %d reports, want 2", len(reports))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		_, err := dbModel.QueryReportsForSummarization(123)

		if err == nil {
			t.Error("QueryReportsForSummarization() expected error for nil DB")
		}
	})
}

func TestQueryReportSummary(t *testing.T) {
	t.Run("summary found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		rows := sqlmock.NewRows([]string{"id", "game_id", "summary"}).
			AddRow(1, 123, "Excellent accessibility overall")

		mock.ExpectQuery("SELECT id, game_id, summary FROM report_summaries WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		summary, err := dbModel.QueryReportSummary(gameID)
		if err != nil {
			t.Errorf("QueryReportSummary() error = %v", err)
		}
		if summary.GameID != 123 {
			t.Errorf("QueryReportSummary() GameID = %d, want 123", summary.GameID)
		}
		if summary.Summary != "Excellent accessibility overall" {
			t.Errorf("QueryReportSummary() Summary = %s", summary.Summary)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("summary not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 999

		mock.ExpectQuery("SELECT id, game_id, summary FROM report_summaries WHERE game_id").
			WithArgs(gameID).
			WillReturnError(sql.ErrNoRows)

		_, err = dbModel.QueryReportSummary(gameID)
		if err == nil {
			t.Error("QueryReportSummary() expected error for not found")
		}
		if err.Error() != "summary not found" {
			t.Errorf("QueryReportSummary() error = %v, want 'summary not found'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		_, err := dbModel.QueryReportSummary(123)

		if err == nil {
			t.Error("QueryReportSummary() expected error for nil DB")
		}
	})
}

func TestQueryUserReports(t *testing.T) {
	t.Run("user reports found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		userID := 42
		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "created_at", "game_id", "game_name", "cover_art", "user_id", "platform", "score", "report"}).
			AddRow(1, now, 123, "Test Game", "cover.jpg", 42, types.PC, 8, "Great!").
			AddRow(2, now, 124, "Another Game", "cover2.jpg", 42, types.Xbox, 7, "Good")

		mock.ExpectQuery("SELECT.*FROM reports r.*LEFT JOIN games g").
			WithArgs(userID).
			WillReturnRows(rows)

		reports, err := dbModel.QueryUserReports(userID)
		if err != nil {
			t.Errorf("QueryUserReports() error = %v", err)
		}
		if len(reports) != 2 {
			t.Errorf("QueryUserReports() returned %d reports, want 2", len(reports))
		}
		if reports[0].GameName != "Test Game" {
			t.Errorf("QueryUserReports() GameName = %s, want 'Test Game'", reports[0].GameName)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		_, err := dbModel.QueryUserReports(42)

		if err == nil {
			t.Error("QueryUserReports() expected error for nil DB")
		}
	})
}
