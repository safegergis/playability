package db

import (
	"errors"
	"fmt"
	"log"
	"playability/types"
)

// QueryFeatureReports retrieves all feature reports from the database
func (m *DatabaseModel) QueryFeatureReports(id int) ([]types.FeatureReport, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		log.Printf("[QueryFeatureReports] Database connection is nil")
		return nil, errors.New("database connection is nil")
	}

	log.Printf("[QueryFeatureReports] Querying feature reports for game ID: %d", id)

	// Query the database for all reports
	query := `SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id = $1`
	rows, err := m.DB.Query(query, id)
	if err != nil {
		log.Printf("[QueryFeatureReports] Error querying feature reports for game ID %d: %v", id, err)
		return nil, fmt.Errorf("error querying feature reports: %w", err)
	}
	defer rows.Close()

	// Iterate through the results and build the accessibility reports slice
	var reports []types.FeatureReport
	for rows.Next() {
		var report types.FeatureReport
		err := rows.Scan(&report.ClosedCaptions, &report.ColorBlind, &report.FullControllerSupport, &report.ControllerRemapping)
		if err != nil {
			log.Printf("[QueryFeatureReports] Error scanning feature report row for game ID %d: %v", id, err)
			return nil, fmt.Errorf("error scanning feature report: %w", err)
		}
		reports = append(reports, report)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("[QueryFeatureReports] Error iterating feature report rows for game ID %d: %v", id, err)
		return nil, fmt.Errorf("error iterating feature reports: %w", err)
	}

	log.Printf("[QueryFeatureReports] Successfully retrieved %d feature reports for game ID %d", len(reports), id)
	return reports, nil
}
