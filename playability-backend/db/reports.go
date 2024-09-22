package db

import (
	"database/sql"
	"errors"
	"log"

	"playability/types"
)

func (m *DatabaseModel) InsertReport(report *types.ReportRow) error {
	if m.DB == nil {
		return errors.New("database connection is nil")
	}
	var existingID int
	checkQuery := `SELECT id FROM reports WHERE game_id = $1 AND user_id = $2`
	err := m.DB.QueryRow(checkQuery, report.GameID, report.UserID).Scan(&existingID)
	if err != sql.ErrNoRows {
		if err != nil {
			log.Println("Error checking report:", err)
			return errors.New("internal server error")
		}
	}
	if existingID != 0 {
		return errors.New("report already exists")
	}
	query := `
	INSERT INTO reports (game_id, user_id, closed_captions, color_blind, full_controller_support, controller_remapping, score, report)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = m.DB.Exec(query, report.GameID, report.UserID, report.ClosedCaptions, report.ColorBlind, report.FullControllerSupport, report.ControllerRemapping, report.Score, report.Report)
	return err
}

func (m *DatabaseModel) QueryReportCards(id string) ([]types.ReportCards, error) {
	if m.DB == nil {
		return nil, errors.New("database connection is nil")
	}
	query := `SELECT id, game_id, user_id, score, report FROM reports WHERE game_id = $1`
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []types.ReportCards
	for rows.Next() {
		var report types.ReportCards
		err := rows.Scan(&report.ID, &report.GameID, &report.UserID, &report.Score, &report.Report)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

func (m *DatabaseModel) QueryAccessibilityReports() ([]types.AccessibilityReport, error) {
	if m.DB == nil {
		return nil, errors.New("database connection is nil")
	}
	query := `SELECT * FROM reports`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []types.AccessibilityReport
	for rows.Next() {
		var report types.AccessibilityReport
		err := rows.Scan(&report.ID, &report.GameID, &report.ClosedCaptions, &report.ColorBlind, &report.FullControllerSupport, &report.ControllerRemapping)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}
