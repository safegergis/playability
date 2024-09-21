package db

import "playability/types"

func (db *DatabaseModel) InsertReport(report *types.ReportRow) error {
	query := `
	INSERT INTO reports (game_id, user_id, closed_captions, color_blind, full_controller_support, controller_remapping, score, report)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := db.DB.Exec(query, report.GameID, report.UserID, report.ClosedCaptions, report.ColorBlind, report.FullControllerSupport, report.ControllerRemapping, report.Score, report.Report)
	return err
}
