package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"playability/types"
)

// InsertReport inserts a new report into the database
// It checks if a report already exists for the given game and user
// If it does, it returns an error
func (m *DatabaseModel) InsertReport(report *types.ReportRow) error {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[InsertReport] Database connection is nil")
		return errors.New("database connection is nil")
	}

	if report == nil {
		log.Printf("[InsertReport] Report is nil")
		return errors.New("report cannot be nil")
	}

	log.Printf("[InsertReport] Attempting to insert report for game ID %d by user ID %d", report.GameID, report.UserID)

	// Check if a report already exists for this game and user
	var existingID int
	checkQuery := `SELECT id FROM reports WHERE game_id = $1 AND user_id = $2`
	err := m.DB.QueryRow(checkQuery, report.GameID, report.UserID).Scan(&existingID)
	if err != sql.ErrNoRows {
		if err != nil {
			log.Printf("[InsertReport] Error checking for existing report (game: %d, user: %d): %v", report.GameID, report.UserID, err)
			return errors.New("internal server error")
		}
	}
	if existingID != 0 {
		log.Printf("[InsertReport] Report already exists for game ID %d by user ID %d (report ID: %d)", report.GameID, report.UserID, existingID)
		return errors.New("report already exists")
	}

	// Insert the new report
	query := `
	INSERT INTO reports (game_id, user_id, closed_captions, color_blind, full_controller_support, controller_remapping, score, report)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`
	var reportID int
	err = m.DB.QueryRow(query, report.GameID, report.UserID, report.ClosedCaptions, report.ColorBlind, report.FullControllerSupport, report.ControllerRemapping, report.Score, report.Report).Scan(&reportID)
	if err != nil {
		log.Printf("[InsertReport] Error inserting report for game ID %d by user ID %d: %v", report.GameID, report.UserID, err)
		return fmt.Errorf("error inserting report: %w", err)
	}

	log.Printf("[InsertReport] Successfully inserted report ID %d for game ID %d by user ID %d (score: %s)", reportID, report.GameID, report.UserID, report.Score)
	return nil
}

// QueryReportCards retrieves report cards for a specific game
func (m *DatabaseModel) QueryReportCards(id int) ([]types.ReportCards, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[QueryReportCards] Database connection is nil")
		return nil, errors.New("database connection is nil")
	}

	log.Printf("[QueryReportCards] Querying report cards for game ID: %d", id)

	// Query the database for reports
	query := `SELECT id, created_at, game_id, user_id, score, report FROM reports WHERE game_id = $1 ORDER BY created_at DESC`
	rows, err := m.DB.Query(query, id)
	if err != nil {
		log.Printf("[QueryReportCards] Error querying reports for game ID %d: %v", id, err)
		return nil, fmt.Errorf("error querying report cards: %w", err)
	}
	defer rows.Close()

	// Iterate through the results and build the report cards slice
	var reports []types.ReportCards
	for rows.Next() {
		var report types.ReportCards
		err := rows.Scan(&report.ID, &report.CreatedAt, &report.GameID, &report.UserID, &report.Score, &report.Report)
		if err != nil {
			log.Printf("[QueryReportCards] Error scanning report row for game ID %d: %v", id, err)
			return nil, fmt.Errorf("error scanning report: %w", err)
		}
		reports = append(reports, report)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("[QueryReportCards] Error iterating report rows for game ID %d: %v", id, err)
		return nil, fmt.Errorf("error iterating reports: %w", err)
	}

	log.Printf("[QueryReportCards] Successfully retrieved %d report cards for game ID %d", len(reports), id)
	return reports, nil
}

func (m *DatabaseModel) QueryAccessibilityScores(id int) ([]int, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[QueryAccessibilityScores] Database connection is nil")
		return nil, errors.New("database connection is nil")
	}

	log.Printf("[QueryAccessibilityScores] Querying accessibility scores for game ID: %d", id)

	// Query the database for the accessibility score
	query := `SELECT score FROM reports WHERE game_id = $1`
	var scores []int
	rows, err := m.DB.Query(query, id)
	if err != nil {
		log.Printf("[QueryAccessibilityScores] Error querying scores for game ID %d: %v", id, err)
		return nil, fmt.Errorf("error querying accessibility scores: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var score int
		err := rows.Scan(&score)
		if err != nil {
			log.Printf("[QueryAccessibilityScores] Error scanning score for game ID %d: %v", id, err)
			return nil, fmt.Errorf("error scanning score: %w", err)
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[QueryAccessibilityScores] Error iterating score rows for game ID %d: %v", id, err)
		return nil, fmt.Errorf("error iterating scores: %w", err)
	}

	log.Printf("[QueryAccessibilityScores] Successfully retrieved %d scores for game ID %d", len(scores), id)
	return scores, nil
}
