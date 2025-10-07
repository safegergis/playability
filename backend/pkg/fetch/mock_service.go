package fetch

import (
	"context"
	"fmt"
)

// MockFetchService is a mock implementation of FetchService for testing
type MockFetchService struct {
	// GetSearchError if set, GetSearch will return this error
	GetSearchError error
	// GetSearchResponse if set, GetSearch will return this response
	GetSearchResponse []byte
	// GetGameError if set, GetGame will return this error
	GetGameError error
	// GetGameResponse if set, GetGame will return this response
	GetGameResponse []byte
	// GetFeaturedGamesError if set, GetFeaturedGames will return this error
	GetFeaturedGamesError error
	// GetFeaturedGamesResponse if set, GetFeaturedGames will return this response
	GetFeaturedGamesResponse []byte
	// CallCount tracks how many times methods were called
	GetSearchCallCount       int
	GetGameCallCount         int
	GetFeaturedGamesCallCount int
	LastSearchTerm           string
	LastGameID               string
}

// GetSearch is a mock implementation that returns pre-configured responses
func (m *MockFetchService) GetSearch(ctx context.Context, searchTerm string) ([]byte, error) {
	m.GetSearchCallCount++
	m.LastSearchTerm = searchTerm

	if m.GetSearchError != nil {
		return nil, m.GetSearchError
	}

	if m.GetSearchResponse != nil {
		return m.GetSearchResponse, nil
	}

	// Default: return empty array
	return []byte(fmt.Sprintf(`[{"id": 1, "name": "%s"}]`, searchTerm)), nil
}

// GetGame is a mock implementation that returns pre-configured responses
func (m *MockFetchService) GetGame(ctx context.Context, gameID string) ([]byte, error) {
	m.GetGameCallCount++
	m.LastGameID = gameID

	if m.GetGameError != nil {
		return nil, m.GetGameError
	}

	if m.GetGameResponse != nil {
		return m.GetGameResponse, nil
	}

	// Default: return mock game data
	return []byte(fmt.Sprintf(`{"id": %s, "name": "Mock Game", "cover_art": "mock.jpg"}`, gameID)), nil
}

// GetFeaturedGames is a mock implementation that returns pre-configured responses
func (m *MockFetchService) GetFeaturedGames(ctx context.Context) ([]byte, error) {
	m.GetFeaturedGamesCallCount++

	if m.GetFeaturedGamesError != nil {
		return nil, m.GetFeaturedGamesError
	}

	if m.GetFeaturedGamesResponse != nil {
		return m.GetFeaturedGamesResponse, nil
	}

	// Default: return mock featured games
	return []byte(`[{"id": 1, "game_id": 123, "name": "Mock Game", "cover_art": "mock.jpg"}]`), nil
}
