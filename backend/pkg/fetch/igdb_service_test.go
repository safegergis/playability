package fetch

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestNewIGDBService(t *testing.T) {
	t.Run("creates service with valid API key", func(t *testing.T) {
		service := NewIGDBService("test-api-key")
		if service == nil {
			t.Error("NewIGDBService() returned nil")
		}
		if service.apiKey != "test-api-key" {
			t.Errorf("NewIGDBService() apiKey = %s, want 'test-api-key'", service.apiKey)
		}
		if service.httpClient == nil {
			t.Error("NewIGDBService() httpClient is nil")
		}
	})

	t.Run("creates service with empty API key", func(t *testing.T) {
		service := NewIGDBService("")
		if service == nil {
			t.Error("NewIGDBService() returned nil")
		}
		// Service should still be created but will fail on API calls
	})
}

func TestIGDBService_GetSearch(t *testing.T) {
	t.Run("empty API key returns error", func(t *testing.T) {
		service := NewIGDBService("")
		ctx := context.Background()

		_, err := service.GetSearch(ctx, "minecraft")

		if err == nil {
			t.Error("GetSearch() should return error when API key is empty")
		}
		if !strings.Contains(err.Error(), "IGDB_ACCESS_TOKEN") {
			t.Errorf("GetSearch() error = %v, want error mentioning IGDB_ACCESS_TOKEN", err)
		}
	})

	t.Run("invalid API key returns error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping API call test in short mode")
		}

		service := NewIGDBService("invalid-key")
		ctx := context.Background()

		_, err := service.GetSearch(ctx, "minecraft")

		if err == nil {
			t.Error("GetSearch() should return error with invalid API key")
		}
	})
}

func TestIGDBService_GetGame(t *testing.T) {
	t.Run("invalid API key returns error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping API call test in short mode")
		}

		service := NewIGDBService("invalid-key")
		ctx := context.Background()

		_, err := service.GetGame(ctx, "123")

		if err == nil {
			t.Error("GetGame() should return error with invalid API key")
		}
	})
}

func TestIGDBService_GetFeaturedGames(t *testing.T) {
	t.Run("empty API key returns error", func(t *testing.T) {
		service := NewIGDBService("")
		ctx := context.Background()

		_, err := service.GetFeaturedGames(ctx)

		if err == nil {
			t.Error("GetFeaturedGames() should return error when API key is empty")
		}
		if !strings.Contains(err.Error(), "IGDB_ACCESS_TOKEN") {
			t.Errorf("GetFeaturedGames() error = %v, want error mentioning IGDB_ACCESS_TOKEN", err)
		}
	})

	t.Run("invalid API key returns error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping API call test in short mode")
		}

		service := NewIGDBService("invalid-key")
		ctx := context.Background()

		_, err := service.GetFeaturedGames(ctx)

		if err == nil {
			t.Error("GetFeaturedGames() should return error with invalid API key")
		}
	})
}

// Integration tests - only run if IGDB_ACCESS_TOKEN is set
func TestIGDBService_GetSearchIntegration(t *testing.T) {
	apiKey := os.Getenv("IGDB_ACCESS_TOKEN")
	if apiKey == "" {
		t.Skip("Skipping integration test - IGDB_ACCESS_TOKEN not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("search for popular game", func(t *testing.T) {
		service := NewIGDBService(apiKey)
		ctx := context.Background()

		result, err := service.GetSearch(ctx, "minecraft")

		if err != nil {
			t.Errorf("GetSearch() unexpected error = %v", err)
		}

		if result == nil || len(result) == 0 {
			t.Error("GetSearch() returned empty result for 'minecraft'")
		}

		// Result should be valid JSON array
		if !strings.HasPrefix(strings.TrimSpace(string(result)), "[") {
			t.Error("GetSearch() result should be a JSON array")
		}
	})
}

func TestIGDBService_GetGameIntegration(t *testing.T) {
	apiKey := os.Getenv("IGDB_ACCESS_TOKEN")
	if apiKey == "" {
		t.Skip("Skipping integration test - IGDB_ACCESS_TOKEN not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("get game details for known game", func(t *testing.T) {
		service := NewIGDBService(apiKey)
		ctx := context.Background()

		// Minecraft game ID on IGDB
		result, err := service.GetGame(ctx, "121")

		if err != nil {
			t.Errorf("GetGame() unexpected error = %v", err)
		}

		if result == nil || len(result) == 0 {
			t.Error("GetGame() returned empty result")
		}

		// Result should contain game data
		resultStr := string(result)
		if !strings.Contains(resultStr, "name") {
			t.Error("GetGame() result should contain 'name' field")
		}
	})
}

func TestIGDBService_GetFeaturedGamesIntegration(t *testing.T) {
	apiKey := os.Getenv("IGDB_ACCESS_TOKEN")
	if apiKey == "" {
		t.Skip("Skipping integration test - IGDB_ACCESS_TOKEN not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("get featured games", func(t *testing.T) {
		service := NewIGDBService(apiKey)
		ctx := context.Background()

		result, err := service.GetFeaturedGames(ctx)

		if err != nil {
			t.Errorf("GetFeaturedGames() unexpected error = %v", err)
		}

		if result == nil || len(result) == 0 {
			t.Error("GetFeaturedGames() returned empty result")
		}

		// Result should be valid JSON array
		resultStr := string(result)
		if !strings.HasPrefix(strings.TrimSpace(resultStr), "[") {
			t.Error("GetFeaturedGames() result should be a JSON array")
		}

		// Should contain game names and cover art
		if !strings.Contains(resultStr, "name") || !strings.Contains(resultStr, "cover_art") {
			t.Error("GetFeaturedGames() result should contain 'name' and 'cover_art' fields")
		}
	})
}
