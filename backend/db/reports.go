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
	//Update user number of reports by one
	query := ` UPDATE users SET num_reports = num_reports + 1 WHERE id = $1`
	_, err = m.DB.Exec(query, report.UserID)
	if err != nil {
		log.Printf("[InsertReport] Error updating num reports by user ID %d: %v", report.UserID, err)
		return fmt.Errorf("error updating num reports by user id: %w", err)
	}
	// Insert the new report
	query = `
	INSERT INTO reports (game_id, user_id, platform, closed_captions, color_blind, full_controller_support, controller_remapping, score, report)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`
	var reportID int
	err = m.DB.QueryRow(query, report.GameID, report.UserID, report.Platform, report.ClosedCaptions, report.ColorBlind, report.FullControllerSupport, report.ControllerRemapping, report.Score, report.Report).Scan(&reportID)
	if err != nil {
		log.Printf("[InsertReport] Error inserting report for game ID %d by user ID %d: %v", report.GameID, report.UserID, err)
		return fmt.Errorf("error inserting report: %w", err)
	}

	log.Printf("[InsertReport] Successfully inserted report ID %d for game ID %d by user ID %d (score: %d)", reportID, report.GameID, report.UserID, report.Score)
	return nil
}
func (m *DatabaseModel) InsertReportSummary(summary *types.ReportSummaryRow) error {

	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[InsertReportSummary] Database connection is nil")
		return errors.New("database connection is nil")
	}

	if summary == nil {
		log.Printf("[InsertReportSummary] summary is nil")
		return errors.New("Summary cannot be nil")
	}

	log.Printf("[InsertReportSummary] Attempting to insert for game ID %d ", summary.GameID)
	query := `
    		INSERT INTO report_summaries (game_id, summary)
    		VALUES ($1, $2 )
    		ON CONFLICT (game_id) DO UPDATE SET
    			summary = EXCLUDED.summary;`
	_, err := m.DB.Exec(query)
	if err != nil {
		log.Printf("[InsertReportSummary] Error upserting report for game ID %d", summary.GameID)
		return errors.New("[InsertReportSummary] Error upserting report ")
	}
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
	query := `SELECT id, created_at, game_id, user_id, platform, score, report FROM reports WHERE game_id = $1 ORDER BY created_at DESC`
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
		err := rows.Scan(&report.ID, &report.CreatedAt, &report.GameID, &report.UserID, &report.Platform, &report.Score, &report.Report)
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

// GetReportCount returns the total number of reports for a specific game
func (m *DatabaseModel) GetReportCount(gameID int) (int, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[GetReportCount] Database connection is nil")
		return 0, errors.New("database connection is nil")
	}

	log.Printf("[GetReportCount] Counting reports for game ID: %d", gameID)

	query := `SELECT COUNT(*) FROM reports WHERE game_id = $1`
	var count int
	err := m.DB.QueryRow(query, gameID).Scan(&count)
	if err != nil {
		log.Printf("[GetReportCount] Error counting reports for game ID %d: %v", gameID, err)
		return 0, fmt.Errorf("error counting reports: %w", err)
	}

	log.Printf("[GetReportCount] Found %d reports for game ID %d", count, gameID)
	return count, nil
}

// QueryReportsForSummarization retrieves all reports for a specific game for AI summarization
func (m *DatabaseModel) QueryReportsForSummarization(gameID int) ([]types.ReportRow, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[QueryReportsForSummarization] Database connection is nil")
		return nil, errors.New("database connection is nil")
	}

	log.Printf("[QueryReportsForSummarization] Querying reports for game ID: %d", gameID)

	query := `SELECT id, game_id, user_id, platform, closed_captions, color_blind, full_controller_support, controller_remapping, score, report FROM reports WHERE game_id = $1 AND report IS NOT NULL AND report != '' ORDER BY created_at DESC`
	rows, err := m.DB.Query(query, gameID)
	if err != nil {
		log.Printf("[QueryReportsForSummarization] Error querying reports for game ID %d: %v", gameID, err)
		return nil, fmt.Errorf("error querying reports for summarization: %w", err)
	}
	defer rows.Close()

	var reports []types.ReportRow
	for rows.Next() {
		var report types.ReportRow
		err := rows.Scan(&report.ID, &report.GameID, &report.UserID, &report.Platform, &report.ClosedCaptions, &report.ColorBlind, &report.FullControllerSupport, &report.ControllerRemapping, &report.Score, &report.Report)
		if err != nil {
			log.Printf("[QueryReportsForSummarization] Error scanning report for game ID %d: %v", gameID, err)
			return nil, fmt.Errorf("error scanning report: %w", err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[QueryReportsForSummarization] Error iterating reports for game ID %d: %v", gameID, err)
		return nil, fmt.Errorf("error iterating reports: %w", err)
	}

	log.Printf("[QueryReportsForSummarization] Successfully retrieved %d reports for game ID %d", len(reports), gameID)
	return reports, nil
}

// QueryReportSummary retrieves the AI-generated summary for a specific game
func (m *DatabaseModel) QueryReportSummary(gameID int) (types.ReportSummaryRow, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[QueryReportSummary] Database connection is nil")
		return types.ReportSummaryRow{}, errors.New("database connection is nil")
	}

	log.Printf("[QueryReportSummary] Querying summary for game ID: %d", gameID)

	query := `SELECT id, game_id, summary FROM report_summaries WHERE game_id = $1`
	var summary types.ReportSummaryRow
	err := m.DB.QueryRow(query, gameID).Scan(&summary.ID, &summary.GameID, &summary.Summary)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[QueryReportSummary] No summary found for game ID: %d", gameID)
			return types.ReportSummaryRow{}, errors.New("summary not found")
		}
		log.Printf("[QueryReportSummary] Error querying summary for game ID %d: %v", gameID, err)
		return types.ReportSummaryRow{}, fmt.Errorf("error querying report summary: %w", err)
	}

	log.Printf("[QueryReportSummary] Successfully retrieved summary for game ID %d", gameID)
	return summary, nil
}

// QueryUserReports retrieves all reports submitted by a specific user with game details
func (m *DatabaseModel) QueryUserReports(userID int) ([]types.ReportCards, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[QueryUserReports] Database connection is nil")
		return nil, errors.New("database connection is nil")
	}

	log.Printf("[QueryUserReports] Querying reports for user ID: %d", userID)

	// Query the database for reports by user, joining with games table to get game details
	query := `
		SELECT
			r.id,
			r.created_at,
			r.game_id,
			COALESCE(g.name, 'Unknown Game') as game_name,
			COALESCE(g.cover_art, '') as cover_art,
			r.user_id,
			r.platform,
			r.score,
			r.report
		FROM reports r
		LEFT JOIN games g ON r.game_id = g.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		log.Printf("[QueryUserReports] Error querying reports for user ID %d: %v", userID, err)
		return nil, fmt.Errorf("error querying user reports: %w", err)
	}
	defer rows.Close()

	// Iterate through the results and build the reports slice
	var reports []types.ReportCards
	for rows.Next() {
		var report types.ReportCards
		err := rows.Scan(&report.ID, &report.CreatedAt, &report.GameID, &report.GameName, &report.CoverArt, &report.UserID, &report.Platform, &report.Score, &report.Report)
		if err != nil {
			log.Printf("[QueryUserReports] Error scanning report row for user ID %d: %v", userID, err)
			return nil, fmt.Errorf("error scanning report: %w", err)
		}
		reports = append(reports, report)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("[QueryUserReports] Error iterating report rows for user ID %d: %v", userID, err)
		return nil, fmt.Errorf("error iterating reports: %w", err)
	}

	log.Printf("[QueryUserReports] Successfully retrieved %d reports for user ID %d", len(reports), userID)
	return reports, nil
}
