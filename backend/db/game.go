package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"playability/types"

	_ "github.com/lib/pq"
)

// InsertGame inserts a new game into the database
func (m DatabaseModel) InsertGame(body []byte) error {
	log.Printf("[InsertGame] Attempting to insert game from %d bytes", len(body))

	// Unmarshal JSON body into Game struct
	var game types.Game
	err := json.Unmarshal(body, &game)
	if err != nil {
		log.Printf("[InsertGame] Error unmarshalling game JSON: %v. Body: %s", err, string(body))
		return fmt.Errorf("error unmarshalling game: %w", err)
	}

	log.Printf("[InsertGame] Inserting game: %s (ID: %d)", game.Name, game.ID)

	// Prepare SQL statement for inserting game
	sqlStatement := `
	INSERT INTO games (id, name, summary, cover_art, platforms, closed_captions, color_blind, full_controller_support, controller_remapping)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		summary = EXCLUDED.summary,
		cover_art = EXCLUDED.cover_art,
		platforms = EXCLUDED.platforms,
		closed_captions = EXCLUDED.closed_captions,
		color_blind = EXCLUDED.color_blind,
		full_controller_support = EXCLUDED.full_controller_support,
		controller_remapping = EXCLUDED.controller_remapping
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
		log.Printf("[InsertGame] Error marshalling platforms for game ID %d: %v", id, err)
		return fmt.Errorf("error marshalling platforms: %w", err)
	}

	// Execute SQL statement
	var returnedID int
	err = m.DB.QueryRow(sqlStatement, id, name, summary, cover_art, platforms, closed_captions, color_blind, full_controller_support, controller_remapping).Scan(&returnedID)
	if err != nil {
		log.Printf("[InsertGame] Error inserting/updating game ID %d (%s): %v", id, name, err)
		return fmt.Errorf("error inserting game: %w", err)
	}

	log.Printf("[InsertGame] Successfully inserted/updated game ID %d (%s)", returnedID, name)
	return nil
}

// QueryGame retrieves a game from the database by its ID
func (m DatabaseModel) QueryGame(id string) ([]byte, error, bool) {
	log.Printf("[QueryGame] Querying game with ID: %s", id)

	// Prepare SQL statement for querying game
	sqlStatement := `
	SELECT id, name, summary, cover_art, platforms, closed_captions, color_blind, full_controller_support, controller_remapping, created_at, updated_at
	FROM games WHERE id = $1`

	var game types.Game
	var platforms string
	var createdAt, updatedAt sql.NullTime

	// Query the database
	row := m.DB.QueryRow(sqlStatement, id)
	err := row.Scan(&game.ID, &game.Name, &game.Summary, &game.CoverArt, &platforms, &game.ClosedCaptions, &game.ColorBlind, &game.FullControllerSupport, &game.ControllerRemapping, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		log.Printf("[QueryGame] No game found with ID: %s", id)
		return nil, nil, false
	} else if err != nil {
		log.Printf("[QueryGame] Error querying game ID %s: %v", id, err)
		return nil, fmt.Errorf("error querying game: %w", err), false
	}

	// Unmarshal platforms JSON string to array
	err = json.Unmarshal([]byte(platforms), &game.Platforms)
	if err != nil {
		log.Printf("[QueryGame] Error unmarshalling platforms for game ID %s: %v. Platforms data: %s", id, err, platforms)
		return nil, fmt.Errorf("error unmarshalling platforms: %w", err), false
	}

	// Marshal game struct to JSON
	body, err := json.Marshal(game)
	if err != nil {
		log.Printf("[QueryGame] Error marshalling game ID %s to JSON: %v", id, err)
		return nil, fmt.Errorf("error marshalling game: %w", err), false
	}

	log.Printf("[QueryGame] Successfully retrieved game: %s (ID: %s)", game.Name, id)
	return body, nil, true
}

// QueryFeaturedGames retrieves cached featured games if they're fresh (< 24 hours old)
func (m *DatabaseModel) QueryFeaturedGames() ([]types.FeaturedGame, bool, error) {
	if m.DB == nil {
		log.Printf("[QueryFeaturedGames] Database connection is nil")
		return nil, false, errors.New("database connection is nil")
	}

	log.Printf("[QueryFeaturedGames] Checking cached featured games")

	// Check the age of the cache by getting the most recent updated_at timestamp
	var lastUpdated sql.NullTime
	err := m.DB.QueryRow("SELECT MAX(updated_at) FROM featured_games").Scan(&lastUpdated)
	if err != nil {
		log.Printf("[QueryFeaturedGames] Error checking cache age: %v", err)
		return nil, false, fmt.Errorf("error checking cache age: %w", err)
	}

	// If no records exist or cache is older than 24 hours, return stale
	if !lastUpdated.Valid || time.Since(lastUpdated.Time) > 24*time.Hour {
		log.Printf("[QueryFeaturedGames] Cache is stale or missing (last updated: %v)", lastUpdated.Time)
		return nil, false, nil
	}

	// Cache is fresh, retrieve the games
	query := `SELECT game_id, name, cover_art FROM featured_games ORDER BY position ASC`
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Printf("[QueryFeaturedGames] Error querying cached featured games: %v", err)
		return nil, false, fmt.Errorf("error querying cached featured games: %w", err)
	}
	defer rows.Close()

	var games []types.FeaturedGame
	for rows.Next() {
		var game types.FeaturedGame
		err := rows.Scan(&game.GameID, &game.Name, &game.CoverArt)
		if err != nil {
			log.Printf("[QueryFeaturedGames] Error scanning featured game row: %v", err)
			return nil, false, fmt.Errorf("error scanning featured game: %w", err)
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[QueryFeaturedGames] Error iterating featured game rows: %v", err)
		return nil, false, fmt.Errorf("error iterating featured games: %w", err)
	}

	log.Printf("[QueryFeaturedGames] Successfully retrieved %d cached featured games", len(games))
	return games, true, nil
}

// UpsertFeaturedGames clears old cache and inserts fresh featured games data
func (m *DatabaseModel) UpsertFeaturedGames(games []types.FeaturedGame) error {
	if m.DB == nil {
		log.Printf("[UpsertFeaturedGames] Database connection is nil")
		return errors.New("database connection is nil")
	}

	log.Printf("[UpsertFeaturedGames] Updating featured games cache with %d games", len(games))

	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		log.Printf("[UpsertFeaturedGames] Error starting transaction: %v", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Clear existing cache
	_, err = tx.Exec("DELETE FROM featured_games")
	if err != nil {
		log.Printf("[UpsertFeaturedGames] Error clearing cache: %v", err)
		return fmt.Errorf("error clearing cache: %w", err)
	}

	// Insert new games
	stmt, err := tx.Prepare("INSERT INTO featured_games (game_id, name, cover_art, position) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Printf("[UpsertFeaturedGames] Error preparing insert statement: %v", err)
		return fmt.Errorf("error preparing insert: %w", err)
	}
	defer stmt.Close()

	for i, game := range games {
		_, err := stmt.Exec(game.GameID, game.Name, game.CoverArt, i+1)
		if err != nil {
			log.Printf("[UpsertFeaturedGames] Error inserting game ID %d: %v", game.GameID, err)
			return fmt.Errorf("error inserting game: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("[UpsertFeaturedGames] Error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Printf("[UpsertFeaturedGames] Successfully cached %d featured games", len(games))
	return nil
}
