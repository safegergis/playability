package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestQueryFeatureReports(t *testing.T) {
	t.Run("successful feature reports query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind", "full_controller_support", "controller_remapping"}).
			AddRow("true", "true", "true", "false").
			AddRow("limited", "true", "false", "true").
			AddRow("false", "limited", "limited", "true")

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		reports, err := dbModel.QueryFeatureReports(gameID)
		if err != nil {
			t.Errorf("QueryFeatureReports() error = %v", err)
		}
		if len(reports) != 3 {
			t.Errorf("QueryFeatureReports() returned %d reports, want 3", len(reports))
		}

		// Verify first report
		if reports[0].ClosedCaptions != "true" {
			t.Errorf("Report 0 ClosedCaptions = %s, want 'true'", reports[0].ClosedCaptions)
		}
		if reports[0].ColorBlind != "true" {
			t.Errorf("Report 0 ColorBlind = %s, want 'true'", reports[0].ColorBlind)
		}
		if reports[0].FullControllerSupport != "true" {
			t.Errorf("Report 0 FullControllerSupport = %s, want 'true'", reports[0].FullControllerSupport)
		}
		if reports[0].ControllerRemapping != "false" {
			t.Errorf("Report 0 ControllerRemapping = %s, want 'false'", reports[0].ControllerRemapping)
		}

		// Verify second report has mixed values
		if reports[1].ClosedCaptions != "limited" {
			t.Errorf("Report 1 ClosedCaptions = %s, want 'limited'", reports[1].ClosedCaptions)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("empty results", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 999

		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind", "full_controller_support", "controller_remapping"})

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		reports, err := dbModel.QueryFeatureReports(gameID)
		if err != nil {
			t.Errorf("QueryFeatureReports() error = %v", err)
		}
		if len(reports) != 0 {
			t.Errorf("QueryFeatureReports() returned %d reports, want 0", len(reports))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnError(errors.New("query error"))

		_, err = dbModel.QueryFeatureReports(gameID)
		if err == nil {
			t.Error("QueryFeatureReports() expected error for query failure")
		}
		if err.Error() != "error querying feature reports: query error" {
			t.Errorf("QueryFeatureReports() error = %v, want 'error querying feature reports: query error'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		// Create row with wrong number of columns to trigger scan error
		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind"}).
			AddRow("true", "true")

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		_, err = dbModel.QueryFeatureReports(gameID)
		if err == nil {
			t.Error("QueryFeatureReports() expected error for scan failure")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("rows iteration error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 123

		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind", "full_controller_support", "controller_remapping"}).
			AddRow("true", "true", "true", "false").
			RowError(0, sql.ErrConnDone)

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		_, err = dbModel.QueryFeatureReports(gameID)
		if err == nil {
			t.Error("QueryFeatureReports() expected error for row iteration failure")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("nil database connection", func(t *testing.T) {
		dbModel := &DatabaseModel{DB: nil}
		_, err := dbModel.QueryFeatureReports(123)

		if err == nil {
			t.Error("QueryFeatureReports() expected error for nil DB")
		}
		if err.Error() != "database connection is nil" {
			t.Errorf("QueryFeatureReports() error = %v, want 'database connection is nil'", err)
		}
	})

	t.Run("all feature value types", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		dbModel := &DatabaseModel{DB: db}
		gameID := 456

		// Test that all three value types (true, limited, false) are properly extracted
		rows := sqlmock.NewRows([]string{"closed_captions", "color_blind", "full_controller_support", "controller_remapping"}).
			AddRow("true", "limited", "false", "true")

		mock.ExpectQuery("SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id").
			WithArgs(gameID).
			WillReturnRows(rows)

		reports, err := dbModel.QueryFeatureReports(gameID)
		if err != nil {
			t.Errorf("QueryFeatureReports() error = %v", err)
		}
		if len(reports) != 1 {
			t.Fatalf("QueryFeatureReports() returned %d reports, want 1", len(reports))
		}

		report := reports[0]
		if report.ClosedCaptions != "true" {
			t.Errorf("ClosedCaptions = %s, want 'true'", report.ClosedCaptions)
		}
		if report.ColorBlind != "limited" {
			t.Errorf("ColorBlind = %s, want 'limited'", report.ColorBlind)
		}
		if report.FullControllerSupport != "false" {
			t.Errorf("FullControllerSupport = %s, want 'false'", report.FullControllerSupport)
		}
		if report.ControllerRemapping != "true" {
			t.Errorf("ControllerRemapping = %s, want 'true'", report.ControllerRemapping)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}
