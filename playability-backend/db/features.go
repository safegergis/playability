package db

import (
	"errors"
	"playability/types"
)

// QueryFeatureReports retrieves all feature reports from the database
func (m *DatabaseModel) QueryFeatureReports(id int) ([]types.FeatureReport, error) {
	// Check if the database connection is valid
	if m.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	// Query the database for all reports
	query := `SELECT closed_captions, color_blind, full_controller_support, controller_remapping FROM reports WHERE game_id = $1`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the results and build the accessibility reports slice
	var reports []types.FeatureReport
	for rows.Next() {
		var report types.FeatureReport
		err := rows.Scan(&report.ID, &report.GameID, &report.ClosedCaptions, &report.ColorBlind, &report.FullControllerSupport, &report.ControllerRemapping)
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
