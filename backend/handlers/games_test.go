package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"playability/db"
	"playability/pkg/fetch"
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

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			FS: &fetch.MockFetchService{},
		}

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

	t.Run("game not in cache - fetches from external API", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		gameID := "123"
		mockGameData := []byte(`{"id": 123, "name": "Fetched Game", "summary": "From API", "cover_art": "api.jpg"}`)

		mockFetch := &fetch.MockFetchService{
			GetGameResponse: mockGameData,
		}

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			FS: mockFetch,
		}

		// QueryGame returns not found (sql.ErrNoRows)
		mock.ExpectQuery("SELECT id, name, summary, cover_art, platforms").
			WithArgs(gameID).
			WillReturnError(sql.ErrNoRows)

		// Expect InsertGame to be called (uses QueryRow with RETURNING)
		mock.ExpectQuery("INSERT INTO games").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

		req := httptest.NewRequest("GET", "/games?id="+gameID, nil)
		w := httptest.NewRecorder()

		env.GetGamesHandler(w, req)

		// Should successfully return the fetched game
		if w.Code != http.StatusOK {
			t.Errorf("GetGamesHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify fetch was called
		if mockFetch.GetGameCallCount != 1 {
			t.Errorf("GetGame called %d times, want 1", mockFetch.GetGameCallCount)
		}

		if mockFetch.LastGameID != gameID {
			t.Errorf("GetGame called with gameID %s, want %s", mockFetch.LastGameID, gameID)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("fetch fails when game not in cache", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		gameID := "456"

		mockFetch := &fetch.MockFetchService{
			GetGameError: sql.ErrConnDone,
		}

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			FS: mockFetch,
		}

		// QueryGame returns not found
		mock.ExpectQuery("SELECT id, name, summary, cover_art, platforms").
			WithArgs(gameID).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest("GET", "/games?id="+gameID, nil)
		w := httptest.NewRecorder()

		env.GetGamesHandler(w, req)

		// Handler should return error when fetch fails
		if w.Code != http.StatusNotFound {
			t.Errorf("GetGamesHandler() status = %d, want %d", w.Code, http.StatusNotFound)
		}

		// Verify fetch was attempted
		if mockFetch.GetGameCallCount != 1 {
			t.Errorf("GetGame called %d times, want 1", mockFetch.GetGameCallCount)
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

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			FS: &fetch.MockFetchService{},
		}

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

	t.Run("stale cache triggers fetch from IGDB", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockFeaturedData := []byte(`[{"id": 1, "game_id": 100, "name": "Fresh Game 1", "cover_art": "fresh1.jpg"}, {"id": 2, "game_id": 101, "name": "Fresh Game 2", "cover_art": "fresh2.jpg"}]`)

		mockFetch := &fetch.MockFetchService{
			GetFeaturedGamesResponse: mockFeaturedData,
		}

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			FS: mockFetch,
		}

		staleTime := time.Now().Add(-25 * time.Hour)

		// Expect cache age check returns stale
		mock.ExpectQuery("SELECT MAX\\(updated_at\\) FROM featured_games").
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(staleTime))

		// Expect UpsertFeaturedGames to be called (transaction with DELETE + prepared INSERT)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM featured_games").
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectPrepare("INSERT INTO featured_games")
		// Expect two Exec calls (one per game)
		mock.ExpectExec("INSERT INTO featured_games").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO featured_games").
			WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("GET", "/featured", nil)
		w := httptest.NewRecorder()

		env.GetFeaturedHandler(w, req)

		// Should successfully return the fetched games
		if w.Code != http.StatusOK {
			t.Errorf("GetFeaturedHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify fetch was called
		if mockFetch.GetFeaturedGamesCallCount != 1 {
			t.Errorf("GetFeaturedGames called %d times, want 1", mockFetch.GetFeaturedGamesCallCount)
		}

		// Verify response contains the fetched data
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

	t.Run("empty cache triggers fetch", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer mockDB.Close()

		mockFeaturedData := []byte(`[{"id": 1, "game_id": 200, "name": "New Game", "cover_art": "new.jpg"}]`)

		mockFetch := &fetch.MockFetchService{
			GetFeaturedGamesResponse: mockFeaturedData,
		}

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			FS: mockFetch,
		}

		// Expect cache age check returns NULL (empty cache)
		mock.ExpectQuery("SELECT MAX\\(updated_at\\) FROM featured_games").
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(nil))

		// Expect UpsertFeaturedGames to be called (transaction with prepared INSERT)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM featured_games").
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectPrepare("INSERT INTO featured_games")
		// Expect one Exec call (one game)
		mock.ExpectExec("INSERT INTO featured_games").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("GET", "/featured", nil)
		w := httptest.NewRecorder()

		env.GetFeaturedHandler(w, req)

		// Should successfully fetch and return games
		if w.Code != http.StatusOK {
			t.Errorf("GetFeaturedHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify fetch was called
		if mockFetch.GetFeaturedGamesCallCount != 1 {
			t.Errorf("GetFeaturedGames called %d times, want 1", mockFetch.GetFeaturedGamesCallCount)
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

		env := &Env{
			DB: db.DatabaseModel{DB: mockDB},
			FS: &fetch.MockFetchService{},
		}

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
	t.Run("successful search returns results", func(t *testing.T) {
		searchTerm := "zelda"
		mockSearchData := []byte(`[{"id": 1, "name": "The Legend of Zelda"}, {"id": 2, "name": "Zelda II"}]`)

		mockFetch := &fetch.MockFetchService{
			GetSearchResponse: mockSearchData,
		}

		env := &Env{
			FS: mockFetch,
		}

		req := httptest.NewRequest("GET", "/search?search="+searchTerm, nil)
		w := httptest.NewRecorder()

		env.GetSearchHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetSearchHandler() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Verify fetch was called with correct search term
		if mockFetch.GetSearchCallCount != 1 {
			t.Errorf("GetSearch called %d times, want 1", mockFetch.GetSearchCallCount)
		}

		if mockFetch.LastSearchTerm != searchTerm {
			t.Errorf("GetSearch called with term %s, want %s", mockFetch.LastSearchTerm, searchTerm)
		}

		// Verify response is valid JSON
		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", w.Header().Get("Content-Type"))
		}
	})

	t.Run("search with empty term", func(t *testing.T) {
		mockFetch := &fetch.MockFetchService{
			GetSearchResponse: []byte(`[]`),
		}

		env := &Env{
			FS: mockFetch,
		}

		req := httptest.NewRequest("GET", "/search?search=", nil)
		w := httptest.NewRecorder()

		env.GetSearchHandler(w, req)

		// Should still call fetch (behavior depends on IGDB API)
		if mockFetch.GetSearchCallCount != 1 {
			t.Errorf("GetSearch called %d times, want 1", mockFetch.GetSearchCallCount)
		}

		if mockFetch.LastSearchTerm != "" {
			t.Errorf("GetSearch called with term %s, want empty string", mockFetch.LastSearchTerm)
		}
	})

	t.Run("fetch error returns error", func(t *testing.T) {
		mockFetch := &fetch.MockFetchService{
			GetSearchError: sql.ErrConnDone,
		}

		env := &Env{
			FS: mockFetch,
		}

		req := httptest.NewRequest("GET", "/search?search=zelda", nil)
		w := httptest.NewRecorder()

		env.GetSearchHandler(w, req)

		// Should return error when fetch fails
		if w.Code != http.StatusNotFound {
			t.Errorf("GetSearchHandler() status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
