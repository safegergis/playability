package fetch

import (
	"context"
)

// FetchService defines the interface for fetching game data from external APIs
type FetchService interface {
	// GetSearch searches for games by search term and returns JSON results
	GetSearch(ctx context.Context, searchTerm string) ([]byte, error)

	// GetGame retrieves detailed information about a specific game by ID
	GetGame(ctx context.Context, gameID string) ([]byte, error)

	// GetFeaturedGames retrieves a list of featured/popular games
	GetFeaturedGames(ctx context.Context) ([]byte, error)
}
