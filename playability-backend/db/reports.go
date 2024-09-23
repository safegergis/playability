package db

import (
	"database/sql"
	"errors"
	"log"

	"playability/types"
)

// InsertReport inserts a new report into the database
// It checks if a report already exists for the given game and user
// If it does, it returns an error
func (m *DatabaseModel) InsertReport(report *types.ReportRow) error {
	// Check if the database connection is valid
	if m.DB == nil {
		return errors.New("database connection is nil")
	}

	// Check if a report already exists for this game and user
	var existingID int
	checkQuery := `SELECT id FROM reports WHERE game_id = $1 AND user_id = $2`
	err := m.DB.QueryRow(checkQuery, report.GameID, report.UserID).Scan(&existingID)
	if err != sql.ErrNoRows {
		if err != nil {
			log.Println("Error checking for existing report:", err)
			return errors.New("internal server error")
		}
	}
	if existingID != 0 {
		return errors.New("report already exists")
	}

	// Insert the new report
	query := `
	INSERT INTO reports (game_id, user_id, closed_captions, color_blind, full_controller_support, controller_remapping, score, report)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = m.DB.Exec(query, report.GameID, report.UserID, report.ClosedCaptions, report.ColorBlind, report.FullControllerSupport, report.ControllerRemapping, report.Score, report.Report)
	if err != nil {
		log.Println("Error inserting report:", err)
		return err
	}

	if err != nil {
		log.Println("Error converting game ID to int:", err)
		return err
	}

	return nil
}

// QueryReportCards retrieves report cards for a specific game
func (m *DatabaseModel) QueryReportCards(id int) ([]types.ReportCards, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	// Query the database for reports
	query := `SELECT id, game_id, user_id, score, report FROM reports WHERE game_id = $1`
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the results and build the report cards slice
	var reports []types.ReportCards
	for rows.Next() {
		var report types.ReportCards
		err := rows.Scan(&report.ID, &report.GameID, &report.UserID, &report.Score, &report.Report)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (m *DatabaseModel) QueryAccessibilityScores(id int) ([]int, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	// Query the database for the accessibility score
	query := `SELECT score FROM reports WHERE game_id = $1`
	var scores []int
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var score int
		err := rows.Scan(&score)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, nil
}
