package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"playability/types"
	"strconv"
)

// IGDBService implements FetchService using IGDB and PCGamingWiki APIs
type IGDBService struct {
	apiKey     string
	httpClient *http.Client
	clientID   string
}

// NewIGDBService creates a new IGDBService with the given API key
func NewIGDBService(apiKey string) *IGDBService {
	if apiKey == "" {
		log.Printf("[IGDBService] Warning: API key is empty - API calls will fail")
	}
	return &IGDBService{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		clientID:   "7bzjkp4ruewaj55ofw7atbugy7p0au",
	}
}

// GetSearch searches for games by search term
func (s *IGDBService) GetSearch(ctx context.Context, searchTerm string) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("IGDB_ACCESS_TOKEN environment variable is not set")
	}

	postBody := fmt.Sprintf("fields id,name;where platforms = (167,169,48,49,6,130) & game_type = (0,8,9) & version_parent = null; search \"%s\"; limit 50;", searchTerm)
	log.Printf("[GetSearch] Searching for: %s", searchTerm)

	body, err := s.makeIgdbRequest(ctx, postBody, "games")
	if err != nil {
		log.Printf("[GetSearch] Error fetching search results for '%s': %v", searchTerm, err)
		return nil, fmt.Errorf("failed to search games: %w", err)
	}

	log.Printf("[GetSearch] Successfully fetched %d bytes for search term '%s'", len(body), searchTerm)
	return body, nil
}

// GetGame retrieves detailed information about a specific game
func (s *IGDBService) GetGame(ctx context.Context, gameID string) ([]byte, error) {
	// Fetch game details from IGDB
	postBody := fmt.Sprintf("fields name,cover,summary,platforms,involved_companies,external_games; where id = %s;", gameID)
	body, err := s.makeIgdbRequest(ctx, postBody, "games")
	if err != nil {
		log.Printf("[GetGame] Error fetching game %s: %v", gameID, err)
		return nil, fmt.Errorf("error making IGDB games request: %w", err)
	}

	log.Printf("[GetGame] Received response for game %s: %d bytes", gameID, len(body))

	// Parse game data
	var games []types.Game
	err = json.Unmarshal(body, &games)
	if err != nil {
		log.Printf("[GetGame] Error unmarshalling game data for ID %s. Response body: %s", gameID, string(body))
		return nil, fmt.Errorf("error unmarshalling IGDB games response: %w", err)
	}

	if len(games) == 0 {
		log.Printf("[GetGame] No game found with ID %s. Response: %s", gameID, string(body))
		return nil, fmt.Errorf("no game found with ID %s", gameID)
	}
	game := games[0]
	log.Printf("[GetGame] Found game: %s (ID: %d)", game.Name, game.ID)

	// Fetch steamID details
	postBody = fmt.Sprintf("fields uid; where game = %d & external_game_source = 1;", game.ID)
	body, err = s.makeIgdbRequest(ctx, postBody, "external_games")
	if err != nil {
		log.Printf("[GetGame] Error fetching external games for game ID %d: %v", game.ID, err)
		return nil, fmt.Errorf("error making IGDB external games request: %w", err)
	}

	var externalGames []types.ExternalGame
	err = json.Unmarshal(body, &externalGames)
	if err != nil {
		log.Printf("[GetGame] Error unmarshalling external games. Response: %s", string(body))
		return nil, fmt.Errorf("error unmarshalling IGDB external games response: %w", err)
	}

	steamID := ""
	steamAvailability := false
	if len(externalGames) == 0 {
		steamAvailability = false
	} else {
		steamAvailability = true
		steamID = externalGames[0].UID
	}

	// Fetch cover art details
	postBody = fmt.Sprintf("fields image_id; where id = %d;", game.Cover)
	body, err = s.makeIgdbRequest(ctx, postBody, "covers")
	if err != nil {
		log.Printf("[GetGame] Error fetching cover art for cover ID %d: %v", game.Cover, err)
		return nil, fmt.Errorf("error making IGDB covers request: %w", err)
	}

	var coverArt []types.CoverArt
	err = json.Unmarshal(body, &coverArt)
	if err != nil {
		log.Printf("[GetGame] Error unmarshalling cover art. Response: %s", string(body))
		return nil, fmt.Errorf("error unmarshalling IGDB covers response: %w", err)
	}

	if len(coverArt) == 0 {
		log.Printf("[GetGame] Warning: No cover art found for cover ID %d", game.Cover)
	}

	// Fetch accessibility information from PCGamingWiki
	closedCaptions := "unknown"
	colorBlind := "unknown"
	fullControllerSupport := "unknown"
	controllerRemapping := "unknown"

	if steamAvailability {
		log.Printf("[GetGame] Fetching accessibility data from PCGamingWiki for Steam ID: %s", steamID)

		// Closed captions
		closedCaptions = s.fetchPCGamingWikiData(ctx, steamID, "Audio", "Audio.Closed_captions")

		// Color blind mode
		colorBlind = s.fetchPCGamingWikiData(ctx, steamID, "Video", "Video.Color_blind")

		// Controller support
		fullControllerSupport, controllerRemapping = s.fetchPCGamingWikiControllerData(ctx, steamID)
	}

	// Set fetched data to game
	if len(coverArt) > 0 {
		game.CoverArt = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", coverArt[0].ImageID)
	}

	game.ClosedCaptions = closedCaptions
	game.ColorBlind = colorBlind
	game.FullControllerSupport = fullControllerSupport
	game.ControllerRemapping = controllerRemapping
	game.SteamAvailability = steamAvailability

	// Marshal game data to JSON
	body, err = json.Marshal(game)
	if err != nil {
		log.Printf("[GetGame] Error marshalling game data for game ID %s: %v", gameID, err)
		return nil, fmt.Errorf("error marshalling game data: %w", err)
	}

	log.Printf("[GetGame] Successfully fetched complete game data for ID %s (%s)", gameID, game.Name)
	return body, nil
}

// GetFeaturedGames retrieves a list of featured/popular games
func (s *IGDBService) GetFeaturedGames(ctx context.Context) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("IGDB_ACCESS_TOKEN environment variable is not set")
	}

	postBody := "fields game_id; sort value desc; limit 10; where popularity_type = 2;"
	log.Printf("[GetFeaturedGames] Fetching featured games")

	body, err := s.makeIgdbRequest(ctx, postBody, "popularity_primitives")
	if err != nil {
		log.Printf("[GetFeaturedGames] Error fetching featured games: %v", err)
		return nil, fmt.Errorf("error making IGDB features request: %w", err)
	}

	var featuredGames []types.FeaturedGame
	err = json.Unmarshal(body, &featuredGames)
	if err != nil {
		log.Printf("[GetFeaturedGames] Error unmarshalling featured games. Response: %s", string(body))
		return nil, fmt.Errorf("error unmarshalling IGDB features response: %w", err)
	}

	log.Printf("[GetFeaturedGames] Found %d featured games", len(featuredGames))

	for i, game := range featuredGames {
		gameDetails, err := s.getFeaturedGameDetails(ctx, strconv.Itoa(game.GameID))
		if err != nil {
			log.Printf("[GetFeaturedGames] Error fetching details for game ID %d: %v", game.GameID, err)
			return nil, fmt.Errorf("error fetching game details for ID %d: %w", game.GameID, err)
		}
		if len(gameDetails) == 0 {
			log.Printf("[GetFeaturedGames] Warning: No details found for game ID %d", game.GameID)
			continue
		}
		featuredGames[i].Name = gameDetails[0].Name
		featuredGames[i].CoverArt = gameDetails[0].CoverArt
	}

	body, err = json.Marshal(featuredGames)
	if err != nil {
		log.Printf("[GetFeaturedGames] Error marshalling featured games: %v", err)
		return nil, fmt.Errorf("error marshalling featured games data: %w", err)
	}

	log.Printf("[GetFeaturedGames] Successfully fetched %d featured games", len(featuredGames))
	return body, nil
}

// getFeaturedGameDetails fetches details for a featured game (private helper)
func (s *IGDBService) getFeaturedGameDetails(ctx context.Context, gameID string) ([]types.FeaturedGameDetailsResponse, error) {
	log.Printf("[getFeaturedGameDetails] Fetching details for game ID %s", gameID)

	postBody := fmt.Sprintf("fields name,cover; where id = %s;", gameID)
	body, err := s.makeIgdbRequest(ctx, postBody, "games")
	if err != nil {
		log.Printf("[getFeaturedGameDetails] Error fetching game %s: %v", gameID, err)
		return nil, fmt.Errorf("error making IGDB games request: %w", err)
	}

	var gameDetails []types.FeaturedGameDetailsResponse
	err = json.Unmarshal(body, &gameDetails)
	if err != nil {
		log.Printf("[getFeaturedGameDetails] Error unmarshalling game details. Response: %s", string(body))
		return nil, fmt.Errorf("error unmarshalling IGDB games response: %w", err)
	}

	if len(gameDetails) == 0 {
		log.Printf("[getFeaturedGameDetails] No game found with ID %s", gameID)
		return nil, fmt.Errorf("no game found with ID %s", gameID)
	}

	postBody = fmt.Sprintf("fields image_id; where id = %d;", gameDetails[0].ImageID)
	body, err = s.makeIgdbRequest(ctx, postBody, "covers")
	if err != nil {
		log.Printf("[getFeaturedGameDetails] Error fetching cover for game %s: %v", gameID, err)
		return nil, fmt.Errorf("error making IGDB covers request: %w", err)
	}

	var coverArt []types.CoverArt
	err = json.Unmarshal(body, &coverArt)
	if err != nil {
		log.Printf("[getFeaturedGameDetails] Error unmarshalling cover art. Response: %s", string(body))
		return nil, fmt.Errorf("error unmarshalling IGDB covers response: %w", err)
	}

	if len(coverArt) == 0 {
		log.Printf("[getFeaturedGameDetails] Warning: No cover art found for game ID %s", gameID)
		gameDetails[0].CoverArt = ""
	} else {
		gameDetails[0].CoverArt = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", coverArt[0].ImageID)
	}

	log.Printf("[getFeaturedGameDetails] Successfully fetched details for game %s (%s)", gameID, gameDetails[0].Name)
	return gameDetails, nil
}

// makeIgdbRequest sends a request to the IGDB API and returns the response (private helper)
func (s *IGDBService) makeIgdbRequest(ctx context.Context, postBody string, endpoint string) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("IGDB access token is empty")
	}

	responseBody := bytes.NewBuffer([]byte(postBody))
	endpnt := "https://api.igdb.com/v4/" + endpoint

	log.Printf("[makeIgdbRequest] Calling IGDB endpoint: %s", endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpnt, responseBody)
	if err != nil {
		log.Printf("[makeIgdbRequest] Error creating request for %s: %v", endpoint, err)
		return nil, fmt.Errorf("error creating IGDB request: %w", err)
	}

	req.Header.Add("Client-ID", s.clientID)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	response, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("[makeIgdbRequest] Error making request to %s: %v", endpoint, err)
		return nil, fmt.Errorf("error making IGDB request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		log.Printf("[makeIgdbRequest] IGDB %s returned status %d. Response: %s", endpoint, response.StatusCode, string(body))
		return nil, fmt.Errorf("IGDB API returned status %d for endpoint %s", response.StatusCode, endpoint)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("[makeIgdbRequest] Error reading response from %s: %v", endpoint, err)
		return nil, fmt.Errorf("error reading IGDB response: %w", err)
	}

	log.Printf("[makeIgdbRequest] Successfully received %d bytes from %s", len(body), endpoint)
	return body, nil
}

// fetchPCGamingWikiData fetches accessibility data from PCGamingWiki (private helper)
func (s *IGDBService) fetchPCGamingWikiData(ctx context.Context, steamID string, table string, field string) string {
	url := fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,%s&fields=%s&join_on=Infobox_game._pageID=%s._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", table, field, table, steamID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("[fetchPCGamingWikiData] Error creating request for %s: %v", field, err)
		return "unknown"
	}

	req.Header.Add("User-Agent", "Playability/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("[fetchPCGamingWikiData] Error making request for %s: %v", field, err)
		return "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[fetchPCGamingWikiData] PCGamingWiki returned status %d for Steam ID %s", resp.StatusCode, steamID)
		return "unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[fetchPCGamingWikiData] Error reading response for %s: %v", field, err)
		return "unknown"
	}

	var response types.PCGamingWikiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("[fetchPCGamingWikiData] Error unmarshalling response for %s. Response: %s", field, string(body))
		return "unknown"
	}

	if len(response.CargoQuery) == 0 {
		log.Printf("[fetchPCGamingWikiData] No data found for %s (Steam ID: %s)", field, steamID)
		return "unknown"
	}

	// Extract the value based on field name
	switch field {
	case "Audio.Closed_captions":
		result := response.CargoQuery[0].Title.ClosedCaptions
		log.Printf("[fetchPCGamingWikiData] Closed captions for Steam ID %s: %s", steamID, result)
		return result
	case "Video.Color_blind":
		result := response.CargoQuery[0].Title.ColorBlind
		log.Printf("[fetchPCGamingWikiData] Color blind mode for Steam ID %s: %s", steamID, result)
		return result
	default:
		return "unknown"
	}
}

// fetchPCGamingWikiControllerData fetches controller support data (private helper)
func (s *IGDBService) fetchPCGamingWikiControllerData(ctx context.Context, steamID string) (string, string) {
	url := fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Input&fields=Input.Full_controller_support,Input.Controller_remapping,&join_on=Infobox_game._pageID=Input._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("[fetchPCGamingWikiControllerData] Error creating request: %v", err)
		return "unknown", "unknown"
	}

	req.Header.Add("User-Agent", "Playability/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("[fetchPCGamingWikiControllerData] Error making request: %v", err)
		return "unknown", "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[fetchPCGamingWikiControllerData] PCGamingWiki returned status %d for Steam ID %s", resp.StatusCode, steamID)
		return "unknown", "unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[fetchPCGamingWikiControllerData] Error reading response: %v", err)
		return "unknown", "unknown"
	}

	var response types.PCGamingWikiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("[fetchPCGamingWikiControllerData] Error unmarshalling response. Response: %s", string(body))
		return "unknown", "unknown"
	}

	if len(response.CargoQuery) == 0 {
		log.Printf("[fetchPCGamingWikiControllerData] No controller data found for Steam ID %s", steamID)
		return "unknown", "unknown"
	}

	fullSupport := response.CargoQuery[0].Title.FullControllerSupport
	remapping := response.CargoQuery[0].Title.ControllerRemapping
	log.Printf("[fetchPCGamingWikiControllerData] Controller support for Steam ID %s - Full: %s, Remapping: %s", steamID, fullSupport, remapping)

	return fullSupport, remapping
}
