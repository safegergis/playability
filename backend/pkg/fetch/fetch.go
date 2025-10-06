package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"playability/types"
	"strconv"
)

func GetSearch(searchTerm string) ([]byte, error) {
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")
	if igdbSecret == "" {
		return nil, fmt.Errorf("IGDB_ACCESS_TOKEN environment variable is not set")
	}

	postBody := fmt.Sprintf("fields id,name;where platforms = (167,169,48,49,6,130) & game_type = (0,8,9) & version_parent = null; search \"%s\"; limit 50;", searchTerm)
	log.Printf("[GetSearch] Searching for: %s", searchTerm)

	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
	if err != nil {
		log.Printf("[GetSearch] Error fetching search results for '%s': %v", searchTerm, err)
		return nil, fmt.Errorf("failed to search games: %w", err)
	}

	log.Printf("[GetSearch] Successfully fetched %d bytes for search term '%s'", len(body), searchTerm)
	return body, nil
}

// getGame retrieves detailed information about a specific game
func GetGame(gameID string) ([]byte, error) {
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")

	// Fetch game details from IGDB
	postBody := fmt.Sprintf("fields name,cover,summary,platforms,involved_companies,external_games; where id = %s;", gameID)
	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
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
	body, err = makeIgdbRequest(postBody, igdbSecret, "external_games")
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
	body, err = makeIgdbRequest(postBody, igdbSecret, "covers")
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
	// Closed captions

	closedCaptions := "unknown"
	colorBlind := "unknown"
	fullControllerSupport := "unknown"
	controllerRemapping := "unknown"
	if steamAvailability {
		// Create a new HTTP request to fetch closed captions information from PCGamingWiki
		log.Printf("[GetGame] Fetching accessibility data from PCGamingWiki for Steam ID: %s", steamID)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Audio&fields=Audio.Closed_captions&join_on=Infobox_game._pageID=Audio._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID), nil)
		if err != nil {
			log.Printf("[GetGame] Error creating PCGamingWiki request for closed captions: %v", err)
			closedCaptions = "unknown"
		} else {
			// Set User-Agent header
			req.Header.Add("User-Agent", "Playability/1.0")
			// Send the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("[GetGame] Error making PCGamingWiki closed captions request: %v", err)
				closedCaptions = "unknown"
			} else if resp.StatusCode != http.StatusOK {
				log.Printf("[GetGame] PCGamingWiki closed captions returned status %d for Steam ID %s", resp.StatusCode, steamID)
				closedCaptions = "unknown"
				resp.Body.Close()
			} else {
				defer resp.Body.Close()
				body, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("[GetGame] Error reading PCGamingWiki closed captions response: %v", err)
					closedCaptions = "unknown"
				} else {
					var response types.PCGamingWikiResponse
					err = json.Unmarshal(body, &response)
					if err != nil {
						log.Printf("[GetGame] Error unmarshalling PCGamingWiki closed captions. Response: %s", string(body))
						closedCaptions = "unknown"
					} else if len(response.CargoQuery) != 0 {
						closedCaptions = response.CargoQuery[0].Title.ClosedCaptions
						log.Printf("[GetGame] Closed captions for Steam ID %s: %s", steamID, closedCaptions)
					} else {
						log.Printf("[GetGame] No closed captions data found for Steam ID %s", steamID)
						closedCaptions = "unknown"
					}
				}
			}
		}

		// Create a new HTTP request to fetch color blind mode information from PCGamingWiki
		req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Video&fields=Video.Color_blind&join_on=Infobox_game._pageID=Video._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID), nil)
		if err != nil {
			log.Printf("[GetGame] Error creating PCGamingWiki color blind request: %v", err)
			colorBlind = "unknown"
		} else {
			// Set User-Agent header
			req.Header.Add("User-Agent", "Playability/1.0")
			// Send the request
			resp, err := http.DefaultClient.Do(req)

			if err != nil {
				log.Printf("[GetGame] Error making PCGamingWiki color blind request: %v", err)
				colorBlind = "unknown"
			} else if resp.StatusCode != http.StatusOK {
				log.Printf("[GetGame] PCGamingWiki color blind returned status %d for Steam ID %s", resp.StatusCode, steamID)
				colorBlind = "unknown"
				resp.Body.Close()
			} else {
				defer resp.Body.Close()
				body, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("[GetGame] Error reading PCGamingWiki color blind response: %v", err)
					colorBlind = "unknown"
				} else {
					var response2 types.PCGamingWikiResponse
					err = json.Unmarshal(body, &response2)
					if err != nil {
						log.Printf("[GetGame] Error unmarshalling PCGamingWiki color blind. Response: %s", string(body))
						colorBlind = "unknown"
					} else if len(response2.CargoQuery) != 0 {
						colorBlind = response2.CargoQuery[0].Title.ColorBlind
						log.Printf("[GetGame] Color blind mode for Steam ID %s: %s", steamID, colorBlind)
					} else {
						log.Printf("[GetGame] No color blind data found for Steam ID %s", steamID)
						colorBlind = "unknown"
					}
				}
			}
		}

		// Create a new HTTP request to fetch controller support and remapping information from PCGamingWiki
		req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Input&fields=Input.Full_controller_support,Input.Controller_remapping,&join_on=Infobox_game._pageID=Input._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID), nil)
		if err != nil {
			log.Printf("[GetGame] Error creating PCGamingWiki controller support request: %v", err)
			fullControllerSupport = "unknown"
			controllerRemapping = "unknown"
		} else {
			// Set User-Agent header
			req.Header.Add("User-Agent", "Playability/1.0")
			// Send the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("[GetGame] Error making PCGamingWiki controller support request: %v", err)
				fullControllerSupport = "unknown"
				controllerRemapping = "unknown"
			} else if resp.StatusCode != http.StatusOK {
				log.Printf("[GetGame] PCGamingWiki controller support returned status %d for Steam ID %s", resp.StatusCode, steamID)
				fullControllerSupport = "unknown"
				controllerRemapping = "unknown"
				resp.Body.Close()
			} else {
				defer resp.Body.Close()
				body, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("[GetGame] Error reading PCGamingWiki controller support response: %v", err)
					fullControllerSupport = "unknown"
					controllerRemapping = "unknown"
				} else {
					var response3 types.PCGamingWikiResponse
					err = json.Unmarshal(body, &response3)
					if err != nil {
						log.Printf("[GetGame] Error unmarshalling PCGamingWiki controller support. Response: %s", string(body))
						fullControllerSupport = "unknown"
						controllerRemapping = "unknown"
					} else if len(response3.CargoQuery) != 0 {
						fullControllerSupport = response3.CargoQuery[0].Title.FullControllerSupport
						controllerRemapping = response3.CargoQuery[0].Title.ControllerRemapping
						log.Printf("[GetGame] Controller support for Steam ID %s - Full: %s, Remapping: %s", steamID, fullControllerSupport, controllerRemapping)
					} else {
						log.Printf("[GetGame] No controller support data found for Steam ID %s", steamID)
						fullControllerSupport = "unknown"
						controllerRemapping = "unknown"
					}
				}
			}
		}
	}

	//set fetched data to game
	game.CoverArt = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", coverArt[0].ImageID)

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

func GetFeaturedGames() ([]byte, error) {
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")
	if igdbSecret == "" {
		return nil, fmt.Errorf("IGDB_ACCESS_TOKEN environment variable is not set")
	}

	postBody := "fields game_id; sort value desc; limit 10; where popularity_type = 2;"
	log.Printf("[GetFeaturedGames] Fetching featured games")

	body, err := makeIgdbRequest(postBody, igdbSecret, "popularity_primitives")
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
		gameDetails, err := getFeaturedGameDetails(strconv.Itoa(game.GameID))
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
func getFeaturedGameDetails(gameID string) ([]types.FeaturedGameDetailsResponse, error) {
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")
	if igdbSecret == "" {
		return nil, fmt.Errorf("IGDB_ACCESS_TOKEN environment variable is not set")
	}

	log.Printf("[getFeaturedGameDetails] Fetching details for game ID %s", gameID)
	postBody := fmt.Sprintf("fields name,cover; where id = %s;", gameID)
	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
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
	body, err = makeIgdbRequest(postBody, igdbSecret, "covers")
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

// makeIgdbRequest sends a request to the IGDB API and returns the response
func makeIgdbRequest(postBody string, secret string, endpoint string) ([]byte, error) {
	if secret == "" {
		return nil, fmt.Errorf("IGDB access token is empty")
	}

	responseBody := bytes.NewBuffer([]byte(postBody))
	endpnt := "https://api.igdb.com/v4/" + endpoint

	log.Printf("[makeIgdbRequest] Calling IGDB endpoint: %s", endpoint)

	req, err := http.NewRequest(http.MethodPost, endpnt, responseBody)
	if err != nil {
		log.Printf("[makeIgdbRequest] Error creating request for %s: %v", endpoint, err)
		return nil, fmt.Errorf("error creating IGDB request: %w", err)
	}

	req.Header.Add("Client-ID", "7bzjkp4ruewaj55ofw7atbugy7p0au")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	response, err := http.DefaultClient.Do(req)
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
