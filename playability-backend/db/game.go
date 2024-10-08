package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"playability/types"

	_ "github.com/lib/pq"
)

// insertGame inserts a new game into the database
func (m DatabaseModel) InsertGame(body []byte) {

	// Unmarshal JSON body into Game struct
	var game types.Game
	err := json.Unmarshal(body, &game)
	if err != nil {
		log.Fatal("Error unmarshalling game:", err)
	}

	// Prepare SQL statement for inserting game
	sqlStatement := `
	INSERT INTO games (id,name, summary, cover_art, platforms, closed_captions, color_blind, full_controller_support, controller_remapping)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`

	// Extract game details
	id := game.ID
	closed_captions := game.ClosedCaptions
	color_blind := game.ColorBlind
	full_controller_support := game.FullControllerSupport
	controller_remapping := game.ControllerRemapping
	name := game.Name
	summary := game.Summary
	cover_art := game.CoverArt

	// Marshal platforms array to JSON string
	platforms, err := json.Marshal(game.Platforms)
	if err != nil {
		log.Fatal("Error marshalling platforms:", err)
	}

	// Execute SQL statement
	_, err = m.DB.Exec(sqlStatement, id, name, summary, cover_art, platforms, closed_captions, color_blind, full_controller_support, controller_remapping)
	if err != nil {
		log.Fatal("Error inserting game:", err)
	}
}

// queryGame retrieves a game from the database by its ID
func (m DatabaseModel) QueryGame(id string) ([]byte, error, bool) {

	// Prepare SQL statement for querying game
	sqlStatement := `
	SELECT * FROM games WHERE id = $1`

	var game types.Game
	var platforms string
	// Query the database
	row := m.DB.QueryRow(sqlStatement, id)
	err := row.Scan(&game.ID, &game.Name, &game.Summary, &game.CoverArt, &platforms, &game.ClosedCaptions, &game.ColorBlind, &game.FullControllerSupport, &game.ControllerRemapping)

	if err == sql.ErrNoRows {
		return nil, nil, false
	} else if err != nil {
		log.Fatal("Error querying game:", err)
	}

	// Unmarshal platforms JSON string to array
	err = json.Unmarshal([]byte(platforms), &game.Platforms)
	if err != nil {
		log.Fatal("Error unmarshalling platforms:", err)
	}

	// Marshal game struct to JSON
	body, err := json.Marshal(game)
	if err != nil {
		log.Fatal("Error marshalling game:", err)
	}
	log.Println("db game:", string(body))
	return body, nil, true
}
