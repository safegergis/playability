package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Database connection constants
const (
	host     = "localhost"
	port     = "5432"
	user     = "safegergis"
	password = "root"
	dbname   = "playability"
)

func main() {
	// Set up the router and define routes
	router := mux.NewRouter()
	router.HandleFunc("/search", searchHandler).Methods("GET")
	router.HandleFunc("/games", gamesHandler).Methods("GET")

	// Start the server on port 8080
	http.ListenAndServe(":8080", router)
}

// enableCors adds CORS headers to allow cross-origin requests
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// searchHandler handles search requests for games
func searchHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// Extract search term from query parameters
	searchTerm := r.URL.Query().Get("search")

	// Call getSearch function (not shown) to perform the search
	body, err := getSearch(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Log and send the search results
	log.Println("Search results:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// gamesHandler handles requests for specific game details
func gamesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// Extract game ID from query parameters
	gameID := r.URL.Query().Get("id")
	// Try to query the game from the database
	body, err, found := queryGame(gameID)
	if !found {
		// If not found in DB, fetch from external source (getGame function not shown)
		body, err = getGame(gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Insert the fetched game into the database
		insertGame(body)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Log and send the game details
	log.Println("game details:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// insertGame inserts a new game into the database
func insertGame(body []byte) {
	// Create database connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Unmarshal JSON body into Game struct
	var game Game
	err = json.Unmarshal(body, &game)
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
	_, err = db.Exec(sqlStatement, id, name, summary, cover_art, platforms, closed_captions, color_blind, full_controller_support, controller_remapping)
	if err != nil {
		log.Fatal("Error inserting game:", err)
	}
}

// queryGame retrieves a game from the database by its ID
func queryGame(id string) ([]byte, error, bool) {
	// Create database connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Prepare SQL statement for querying game
	sqlStatement := `
	SELECT * FROM games WHERE id = $1`

	var game Game
	var platforms string
	// Query the database
	row := db.QueryRow(sqlStatement, id)
	err = row.Scan(&game.ID, &game.Name, &game.Summary, &game.CoverArt, &platforms, &game.ClosedCaptions, &game.ColorBlind, &game.FullControllerSupport, &game.ControllerRemapping)

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

	log.Println("game:", game)
	// Marshal game struct to JSON
	body, err := json.Marshal(game)
	if err != nil {
		log.Fatal("Error marshalling game:", err)
	}
	log.Println("db game:", string(body))
	return body, nil, true
}
