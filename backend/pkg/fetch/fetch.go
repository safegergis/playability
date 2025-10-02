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

	"github.com/joho/godotenv"
)

func GetSearch(searchTerm string) ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")

	postBody := fmt.Sprintf("fields id,name;where platforms = (167,168,48,49,6,130) & category = (0,8,9); search \"%s\"; limit 50;", searchTerm)

	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
	if err != nil {
		log.Fatal(err)
	}

	return body, nil
}

// getGame retrieves detailed information about a specific game
func GetGame(gameID string) ([]byte, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")

	// Fetch game details from IGDB
	postBody := fmt.Sprintf("fields name,cover,summary,platforms,involved_companies,external_games; where id = %s;", gameID)
	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
	if err != nil {
		return nil, fmt.Errorf("error making IGDB games request: %w", err)
	}

	// Parse game data
	var games []types.Game
	err = json.Unmarshal(body, &games)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling IGDB games response: %w", err)
	}

	if len(games) == 0 {
		return nil, fmt.Errorf("no game found with ID %s", gameID)
	}
	game := games[0]
	// Fetch steamID details
	postBody = fmt.Sprintf("fields uid; where game = %d & category = 1;", game.ID)
	body, err = makeIgdbRequest(postBody, igdbSecret, "external_games")
	if err != nil {
		log.Fatal("Error making IGDB external games request:", err)
	}

	var externalGames []types.ExternalGame
	err = json.Unmarshal(body, &externalGames)
	if err != nil {
		log.Fatal("Error unmarshalling IGDB external games response:", err)
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
		log.Fatal("Error making IGDB covers request:", err)
	}

	var coverArt []types.CoverArt
	err = json.Unmarshal(body, &coverArt)
	if err != nil {
		log.Fatal("Error unmarshalling IGDB covers response:", err)
	}

	// Fetch accessibility information from PCGamingWiki
	// Closed captions

	closedCaptions := "unknown"
	colorBlind := "unknown"
	fullControllerSupport := "unknown"
	controllerRemapping := "unknown"
	if steamAvailability {
		// Create a new HTTP request to fetch closed captions information from PCGamingWiki
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Audio&fields=Audio.Closed_captions&join_on=Infobox_game._pageID=Audio._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID), nil)
		if err != nil {
			log.Fatal("Error creating PC Gaming Wiki closed captions request:", err)
		}
		// Set User-Agent header
		req.Header.Add("User-Agent", "Playability/1.0")
		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if resp.StatusCode != http.StatusOK || err != nil {
			log.Println("Error making PC Gaming Wiki closed captions request:", err, resp.StatusCode)
			closedCaptions = "unknown"
		} else {
			// Ensure the response body is closed after we're done
			defer resp.Body.Close()
			// Read the response body
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading PC Gaming Wiki closed captions:", err)
			}
			// Parse the JSON response
			var response types.PCGamingWikiResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.Fatal("Error unmarshalling PC Gaming Wiki closed captions:", err)
			}
			// Extract closed captions information if available
			if len(response.CargoQuery) != 0 {
				closedCaptions = response.CargoQuery[0].Title.ClosedCaptions
				if err != nil {
					log.Fatal("Error parsing closed captions:", err)
				}
			} else {
				closedCaptions = "unknown"
			}
		}

		// Create a new HTTP request to fetch color blind mode information from PCGamingWiki
		req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Video&fields=Video.Color_blind&join_on=Infobox_game._pageID=Video._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID), nil)
		if err != nil {
			log.Println("Error creating PC Gaming Wiki color blind mode request:", err)
		}
		// Set User-Agent header
		req.Header.Add("User-Agent", "Playability/1.0")
		// Send the request
		resp, err = http.DefaultClient.Do(req)

		if resp.StatusCode != http.StatusOK || err != nil {
			log.Println("Error making PC Gaming Wiki color blind mode request:", err, resp.StatusCode)
			colorBlind = "unknown"
		} else {
			// Ensure the response body is closed after we're done
			defer resp.Body.Close()
			// Read the response body
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading PC Gaming Wiki color blind mode:", err)
			}

			// Parse the JSON response
			var response2 types.PCGamingWikiResponse
			err = json.Unmarshal(body, &response2)
			if err != nil {
				log.Fatal("Error unmarshalling PC Gaming Wiki color blind mode:", err)
			}
			// Extract color blind mode information if available
			if len(response2.CargoQuery) != 0 {
				colorBlind = response2.CargoQuery[0].Title.ColorBlind
				if err != nil {
					log.Fatal("Error parsing color blind mode:", err)
				}
			} else {
				colorBlind = "unknown"
			}
		}

		// Create a new HTTP request to fetch controller support and remapping information from PCGamingWiki
		req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Input&fields=Input.Full_controller_support,Input.Controller_remapping,&join_on=Infobox_game._pageID=Input._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID), nil)
		if err != nil {
			log.Fatal("Error creating PC Gaming Wiki controller support and remapping request:", err)
		}
		// Set User-Agent header
		req.Header.Add("User-Agent", "Playability/1.0")
		// Send the request
		resp, err = http.DefaultClient.Do(req)
		if resp.StatusCode != http.StatusOK || err != nil {
			log.Println("Error making PC Gaming Wiki controller support and remapping request:", err)
			fullControllerSupport = "unknown"
			controllerRemapping = "unknown"
		} else {
			// Ensure the response body is closed after we're done
			defer resp.Body.Close()
			// Read the response body
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading PC Gaming Wiki controller support and remapping:", err)
			}
			// Parse the JSON response
			var response3 types.PCGamingWikiResponse
			err = json.Unmarshal(body, &response3)
			if err != nil {
				log.Fatal("Error unmarshalling PC Gaming Wiki controller support and remapping:", err)
			}
			// Extract controller support and remapping information if available
			if len(response3.CargoQuery) != 0 {
				fullControllerSupport = response3.CargoQuery[0].Title.FullControllerSupport
				if err != nil {
					log.Fatal("Error parsing full controller support:", err)
				}
				controllerRemapping = response3.CargoQuery[0].Title.ControllerRemapping
				if err != nil {
					log.Fatal("Error parsing controller remapping:", err)
				}
			} else {
				fullControllerSupport = "unknown"
				controllerRemapping = "unknown"
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
		log.Fatal("Error marshalling game data:", err)
	}

	return body, nil

}

func GetFeaturedGames() ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")

	postBody := "fields game_id; sort value desc; limit 10; where popularity_type = 3;"

	body, err := makeIgdbRequest(postBody, igdbSecret, "popularity_primitives")
	if err != nil {
		log.Fatal("Error making IGDB features request:", err)
	}

	var featuredGames []types.FeaturedGame
	err = json.Unmarshal(body, &featuredGames)
	if err != nil {
		log.Fatal("Error unmarshalling IGDB features response:", err)
	}

	for i, game := range featuredGames {
		gameDetails, err := getFeaturedGameDetails(strconv.Itoa(game.GameID))
		if err != nil {
			log.Fatal("Error fetching game details:", err)
		}
		featuredGames[i].Name = gameDetails[0].Name
		featuredGames[i].CoverArt = gameDetails[0].CoverArt
	}
	body, err = json.Marshal(featuredGames)
	if err != nil {
		log.Fatal("Error marshalling featured games data:", err)
	}
	return body, nil
}
func getFeaturedGameDetails(gameID string) ([]types.FeaturedGameDetailsResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")
	postBody := fmt.Sprintf("fields name,cover; where id = %s;", gameID)
	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
	if err != nil {
		log.Fatal("Error making IGDB games request:", err)
	}

	var gameDetails []types.FeaturedGameDetailsResponse
	err = json.Unmarshal(body, &gameDetails)
	if err != nil {
		log.Fatal("Error unmarshalling IGDB games response:", err)
	}
	postBody = fmt.Sprintf("fields image_id; where id = %d;", gameDetails[0].ImageID)
	body, err = makeIgdbRequest(postBody, igdbSecret, "covers")
	if err != nil {
		log.Fatal("Error making IGDB covers request:", err)
	}

	var coverArt []types.CoverArt
	err = json.Unmarshal(body, &coverArt)
	if err != nil {
		log.Fatal("Error unmarshalling IGDB covers response:", err)
	}
	gameDetails[0].CoverArt = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", coverArt[0].ImageID)
	return gameDetails, nil
}

// makeIgdbRequest sends a request to the IGDB API and returns the response
func makeIgdbRequest(postBody string, secret string, endpoint string) ([]byte, error) {
	responseBody := bytes.NewBuffer([]byte(postBody))

	endpnt := "https://api.igdb.com/v4/" + endpoint

	resp, err := http.NewRequest(http.MethodPost, endpnt, responseBody)
	resp.Header.Add("Client-ID", "7bzjkp4ruewaj55ofw7atbugy7p0au")
	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
