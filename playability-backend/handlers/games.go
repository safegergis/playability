package handlers

import (
	"log"
	"net/http"

	"playability/pkg/fetch"
)

// searchHandler handles search requests for games
func GetSearchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract search term from query parameters
	searchTerm := r.URL.Query().Get("search")

	// Call getSearch function (not shown) to perform the search
	body, err := fetch.GetSearch(searchTerm)
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
func (env *Env) GetGamesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from query parameters
	gameID := r.URL.Query().Get("id")
	// Try to query the game from the database
	body, err, found := env.DB.QueryGame(gameID)
	if !found {
		// If not found in DB, fetch from external source (getGame function not shown)
		body, err = fetch.GetGame(gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Insert the fetched game into the database
		env.DB.InsertGame(body)
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
