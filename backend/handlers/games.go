package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"playability/pkg/fetch"
	"playability/types"
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
func (env *Env) GetFeaturedHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GetFeaturedHandler] Request for featured games from %s", r.RemoteAddr)

	// Try to get cached featured games
	cachedGames, isFresh, err := env.DB.QueryFeaturedGames()
	if err != nil {
		log.Printf("[GetFeaturedHandler] Error checking cache: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// If cache is fresh, return it
	if isFresh && len(cachedGames) > 0 {
		log.Printf("[GetFeaturedHandler] Returning %d cached featured games", len(cachedGames))
		body, err := json.Marshal(cachedGames)
		if err != nil {
			log.Printf("[GetFeaturedHandler] Error marshalling cached games: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
		return
	}

	// Cache is stale or empty, fetch from IGDB
	log.Printf("[GetFeaturedHandler] Cache stale/empty, fetching from IGDB")
	body, err := fetch.GetFeaturedGames()
	if err != nil {
		log.Printf("[GetFeaturedHandler] Error fetching from IGDB: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Parse the fetched games
	var freshGames []types.FeaturedGame
	err = json.Unmarshal(body, &freshGames)
	if err != nil {
		log.Printf("[GetFeaturedHandler] Error unmarshalling IGDB response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Update the cache
	err = env.DB.UpsertFeaturedGames(freshGames)
	if err != nil {
		log.Printf("[GetFeaturedHandler] Warning: Failed to cache games: %v", err)
		// Don't fail the request, just log the warning
	} else {
		log.Printf("[GetFeaturedHandler] Successfully cached %d featured games", len(freshGames))
	}

	// Return the freshly fetched games
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
