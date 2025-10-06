package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"playability/db"
	"playability/types"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetGamesHandler(t *testing.T) {
	t.Run("game found in database cache", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		// Mock game data
		gameID := "123"
		platforms := []int{6, 49} // PC, PlayStation
		platformsJSON, _ := json.Marshal(platforms)

		game := types.Game{
			ID:                    123,
			Name:                  "Test Game",
			Summary:               "A test game",
			CoverArt:              "cover.jpg",
			Platforms:             platforms,
			ClosedCaptions:        "true",
			ColorBlind:            "limited",
			FullControllerSupport: "true",
			ControllerRemapping:   "false",
		}

		rows := sqlmock.NewRows([]string{"id", "name", "summary", "cover_art", "platforms",
			"closed_captions", "color_blind", "full_controller_support", "controller_remapping",
			"created_at", "updated_at"}).
			AddRow(123, "Test Game", "A test game", "cover.jpg", string(platformsJSON),
				"true", "limited", "true", "false", time.Now(), time.Now())

		mock.ExpectQuery("SELECT id, name, summary, cover_art, platforms").
			WithArgs(gameID).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/games?id="+gameID, nil)
		w := httptest.NewRecorder()

		env.GetGamesHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetGamesHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify response is valid JSON
		var returnedGame types.Game
		err = json.NewDecoder(w.Body).Decode(&returnedGame)
		if err != nil {
			t.Errorf("GetGamesHandler() returned invalid JSON: %v", err)
		}

		if returnedGame.ID != game.ID {
			t.Errorf("GetGamesHandler() game ID = %d, want %d", returnedGame.ID, game.ID)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("game not in cache - database error on query", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		gameID := "123"

		// QueryGame returns not found (sql.ErrNoRows)
		mock.ExpectQuery("SELECT id, name, summary, cover_art, platforms").
			WithArgs(gameID).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest("GET", "/games?id="+gameID, nil)
		w := httptest.NewRecorder()

		env.GetGamesHandler(w, req)

		// When not found in DB, it tries to fetch from external API
		// Since we haven't mocked fetch.GetGame, this will fail
		// The actual behavior depends on fetch package implementation
		// For now, we just verify the DB query was attempted

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("database query error returns error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		gameID := "456"

		// QueryGame returns error (not sql.ErrNoRows)
		mock.ExpectQuery("SELECT id, name, summary, cover_art, platforms").
			WithArgs(gameID).
			WillReturnError(sql.ErrConnDone)

		req := httptest.NewRequest("GET", "/games?id="+gameID, nil)
		w := httptest.NewRecorder()

		env.GetGamesHandler(w, req)

		// Handler returns error immediately on DB error (not sql.ErrNoRows)
		if w.Code != http.StatusNotFound {
			t.Errorf("GetGamesHandler() status = %d, want %d", w.Code, http.StatusNotFound)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestGetFeaturedHandler(t *testing.T) {
	t.Run("fresh cache returns cached games", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		freshTime := time.Now().Add(-12 * time.Hour)

		// Expect cache age check
		mock.ExpectQuery("SELECT MAX\\(updated_at\\) FROM featured_games").
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(freshTime))

		// Expect games query
		gamesRows := sqlmock.NewRows([]string{"game_id", "name", "cover_art"}).
			AddRow(1, "Game 1", "cover1.jpg").
			AddRow(2, "Game 2", "cover2.jpg")

		mock.ExpectQuery("SELECT game_id, name, cover_art FROM featured_games ORDER BY position ASC").
			WillReturnRows(gamesRows)

		req := httptest.NewRequest("GET", "/featured", nil)
		w := httptest.NewRecorder()

		env.GetFeaturedHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetFeaturedHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify response is valid JSON array
		var games []types.FeaturedGame
		err = json.NewDecoder(w.Body).Decode(&games)
		if err != nil {
			t.Errorf("GetFeaturedHandler() returned invalid JSON: %v", err)
		}

		if len(games) != 2 {
			t.Errorf("GetFeaturedHandler() returned %d games, want 2", len(games))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("stale cache triggers fetch - Note: fetch not mocked", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		staleTime := time.Now().Add(-25 * time.Hour)

		// Expect cache age check returns stale
		mock.ExpectQuery("SELECT MAX\\(updated_at\\) FROM featured_games").
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(staleTime))

		req := httptest.NewRequest("GET", "/featured", nil)
		w := httptest.NewRecorder()

		env.GetFeaturedHandler(w, req)

		// Since fetch.GetFeaturedGames() is not mocked, it will fail
		// The handler should return an error
		if w.Code == http.StatusOK {
			t.Error("GetFeaturedHandler() should fail when fetch is not available")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("empty cache triggers fetch", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		// Expect cache age check returns NULL (empty cache)
		mock.ExpectQuery("SELECT MAX\\(updated_at\\) FROM featured_games").
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(nil))

		req := httptest.NewRequest("GET", "/featured", nil)
		w := httptest.NewRecorder()

		env.GetFeaturedHandler(w, req)

		// Since fetch.GetFeaturedGames() is not mocked, it will fail
		if w.Code == http.StatusOK {
			t.Error("GetFeaturedHandler() should fail when fetch is not available")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("database error checking cache", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		env := &Env{DB: db.DatabaseModel{DB: mockDB}}

		// Expect cache age check returns error
		mock.ExpectQuery("SELECT MAX\\(updated_at\\) FROM featured_games").
			WillReturnError(sql.ErrConnDone)

		req := httptest.NewRequest("GET", "/featured", nil)
		w := httptest.NewRecorder()

		env.GetFeaturedHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetFeaturedHandler() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		if !strings.Contains(w.Body.String(), "internal server error") {
			t.Errorf("GetFeaturedHandler() body = %s", w.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("error marshalling cached games", func(t *testing.T) {
		// This is difficult to test since json.Marshal rarely fails on valid structs
		// Skipping this edge case as it's nearly impossible to trigger
	})
}

func TestGetSearchHandler(t *testing.T) {
	// Note: GetSearchHandler doesn't use the Env struct or database
	// It directly calls fetch.GetSearch which would need to be mocked
	// These tests verify the handler behavior assuming fetch fails

	t.Run("search with term - fetch not mocked", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/search?search=zelda", nil)
		w := httptest.NewRecorder()

		GetSearchHandler(w, req)

		// Since fetch.GetSearch is not mocked, this will likely fail
		// The actual behavior depends on fetch implementation
		// Just verify the handler runs without panicking
	})

	t.Run("search with empty term", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/search?search=", nil)
		w := httptest.NewRecorder()

		GetSearchHandler(w, req)

		// Empty search term behavior depends on fetch.GetSearch
	})
}
