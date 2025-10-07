package fetch

import (
	"os"
	"strings"
	"testing"
)

// Note: These tests require IGDB_ACCESS_TOKEN to be set and will make real API calls
// For true unit testing, the HTTP client should be dependency-injected
// Current implementation uses environment variables and http.DefaultClient

func TestGetSearch(t *testing.T) {
	t.Run("missing API token", func(t *testing.T) {
		originalToken := os.Getenv("IGDB_ACCESS_TOKEN")
		os.Unsetenv("IGDB_ACCESS_TOKEN")
		defer os.Setenv("IGDB_ACCESS_TOKEN", originalToken)

		_, err := GetSearch("minecraft")

		if err == nil {
			t.Error("GetSearch() should return error when API token is missing")
		}
		if !strings.Contains(err.Error(), "IGDB_ACCESS_TOKEN environment variable is not set") {
			t.Errorf("GetSearch() error = %v, want 'IGDB_ACCESS_TOKEN environment variable is not set'", err)
		}
	})

	t.Run("empty search term - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("valid search results - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")

		// Example of what the test would look like with proper DI:
		// mockClient := &MockHTTPClient{
		//     Response: &http.Response{
		//         StatusCode: 200,
		//         Body: io.NopCloser(strings.NewReader(`[{"id": 1, "name": "Minecraft"}]`)),
		//     },
		// }
		//
		// result, err := GetSearchWithClient("minecraft", mockClient, "test-token")
		// if err != nil {
		//     t.Errorf("GetSearch() unexpected error = %v", err)
		// }
		// ...
	})

	t.Run("IGDB API error - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})
}

func TestGetGame(t *testing.T) {
	t.Run("missing API token", func(t *testing.T) {
		originalToken := os.Getenv("IGDB_ACCESS_TOKEN")
		os.Unsetenv("IGDB_ACCESS_TOKEN")
		defer os.Setenv("IGDB_ACCESS_TOKEN", originalToken)

		_, err := GetGame("123")

		if err == nil {
			t.Error("GetGame() should return error when API token is missing")
		}
	})

	t.Run("valid game with Steam ID - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("game without Steam ID - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("game not found - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("cover art fetching - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("PCGamingWiki accessibility features - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})
}

func TestGetFeaturedGames(t *testing.T) {
	t.Run("missing API token", func(t *testing.T) {
		originalToken := os.Getenv("IGDB_ACCESS_TOKEN")
		os.Unsetenv("IGDB_ACCESS_TOKEN")
		defer os.Setenv("IGDB_ACCESS_TOKEN", originalToken)

		_, err := GetFeaturedGames()

		if err == nil {
			t.Error("GetFeaturedGames() should return error when API token is missing")
		}
		if !strings.Contains(err.Error(), "IGDB_ACCESS_TOKEN environment variable is not set") {
			t.Errorf("GetFeaturedGames() error = %v, want 'IGDB_ACCESS_TOKEN environment variable is not set'", err)
		}
	})

	t.Run("popularity primitives fetch - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("game details enrichment - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("marshaling error handling - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})
}

func TestGetFeaturedGameDetails(t *testing.T) {
	t.Run("missing API token", func(t *testing.T) {
		originalToken := os.Getenv("IGDB_ACCESS_TOKEN")
		os.Unsetenv("IGDB_ACCESS_TOKEN")
		defer os.Setenv("IGDB_ACCESS_TOKEN", originalToken)

		_, err := getFeaturedGameDetails("123")

		if err == nil {
			t.Error("getFeaturedGameDetails() should return error when API token is missing")
		}
		if !strings.Contains(err.Error(), "IGDB_ACCESS_TOKEN environment variable is not set") {
			t.Errorf("getFeaturedGameDetails() error = %v, want 'IGDB_ACCESS_TOKEN environment variable is not set'", err)
		}
	})

	t.Run("valid game fetch - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("game not found - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("cover art fetching - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})
}

func TestMakeIgdbRequest(t *testing.T) {
	t.Run("empty API token", func(t *testing.T) {
		_, err := makeIgdbRequest("test body", "", "games")

		if err == nil {
			t.Error("makeIgdbRequest() should return error when token is empty")
		}
		if !strings.Contains(err.Error(), "IGDB access token is empty") {
			t.Errorf("makeIgdbRequest() error = %v, want 'IGDB access token is empty'", err)
		}
	})

	t.Run("successful request with headers - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")

		// With proper DI, we would test:
		// - Client-ID header is set correctly
		// - Authorization Bearer token is set
		// - POST body is sent correctly
		// - Response status codes are handled
		// - Response body is read correctly
	})

	t.Run("API error response handling - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})

	t.Run("non-200 status codes - requires API mocking", func(t *testing.T) {
		t.Skip("Requires HTTP client mocking - client should be dependency injected")
	})
}

// Integration tests - only run if IGDB_ACCESS_TOKEN is set
func TestGetSearchIntegration(t *testing.T) {
	if os.Getenv("IGDB_ACCESS_TOKEN") == "" {
		t.Skip("Skipping integration test - IGDB_ACCESS_TOKEN not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("search for popular game", func(t *testing.T) {
		result, err := GetSearch("minecraft")

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

func TestGetGameIntegration(t *testing.T) {
	if os.Getenv("IGDB_ACCESS_TOKEN") == "" {
		t.Skip("Skipping integration test - IGDB_ACCESS_TOKEN not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("get game details for known game", func(t *testing.T) {
		// Minecraft game ID on IGDB
		result, err := GetGame("121")

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

func TestGetFeaturedGamesIntegration(t *testing.T) {
	if os.Getenv("IGDB_ACCESS_TOKEN") == "" {
		t.Skip("Skipping integration test - IGDB_ACCESS_TOKEN not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("get featured games", func(t *testing.T) {
		result, err := GetFeaturedGames()

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
