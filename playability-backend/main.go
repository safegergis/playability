package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"bytes"
)

// Game represents the structure of a game from the IGDB API
type Game struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Cover                 int    `json:"cover"`
	Summary               string `json:"summary"`
	Platforms             []int  `json:"platforms"`
	InvolvedCompanies     []int  `json:"involved_companies"`
	ExternalGames         []int  `json:"external_games"`
	CoverArt              string `json:"cover_art"`
	ClosedCaptions        bool   `json:"closed_captions"`
	ColorBlind            bool   `json:"color_blind"`
	FullControllerSupport bool   `json:"full_controller_support"`
	ControllerRemapping   bool   `json:"controller_remapping"`
	SteamAvailability     bool   `json:"steam_availability"`
}

// CoverArt represents the structure of cover art data from the IGDB API
type CoverArt struct {
	ID      int    `json:"id"`
	ImageID string `json:"image_id"`
}

type ExternalGame struct {
	ID  int    `json:"id"`
	UID string `json:"uid"`
}

// PCGamingWikiResponse represents the structure of responses from PCGamingWiki API
type PCGamingWikiResponse struct {
	CargoQuery []struct {
		Title struct {
			ClosedCaptions        string `json:"Closed captions,omitempty"`
			ColorBlind            string `json:"Color blind,omitempty"`
			FullControllerSupport string `json:"Full controller support,omitempty"`
			ControllerRemapping   string `json:"Controller remapping,omitempty"`
		} `json:"title"`
	} `json:"cargoquery"`
}

func main() {
	// Set up the router and define routes
	router := mux.NewRouter()
	router.HandleFunc("/search", searchHandler).Methods("GET")
	router.HandleFunc("/games", gamesHandler).Methods("GET")

	// Start the server
	http.ListenAndServe(":8080", router)
}

// enableCors adds CORS headers to allow cross-origin requests
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// searchHandler handles search requests for games
func searchHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	searchTerm := r.URL.Query().Get("search")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")

	postBody := fmt.Sprintf("fields id,name;where platforms = (167,168,48,49,6,130); search \"%s\"; limit 50;", searchTerm)

	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Search results:", string(body))
	w.Write([]byte(body))
	w.Header().Set("Content-Type", "application/json")

}

// gamesHandler handles requests for specific game details
func gamesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	gameID := r.URL.Query().Get("id")
	body, err := getGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Println("game details:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// getGame retrieves detailed information about a specific game
func getGame(gameID string) ([]byte, error) {
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
	var game []Game
	err = json.Unmarshal(body, &game)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling IGDB games response: %w", err)
	}

	if len(game) == 0 {
		return nil, fmt.Errorf("no game found with ID %s", gameID)
	}

	// Fetch steamID details
	postBody = fmt.Sprintf("fields uid; where game = %d & category = 1;", game[0].ID)
	body, err = makeIgdbRequest(postBody, igdbSecret, "external_games")
	if err != nil {
		log.Fatal("Error making IGDB external games request:", err)
	}

	var externalGames []ExternalGame
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

	log.Println("Steam ID:", steamID)

	// Fetch cover art details
	postBody = fmt.Sprintf("fields image_id; where id = %d;", game[0].Cover)
	body, err = makeIgdbRequest(postBody, igdbSecret, "covers")
	if err != nil {
		log.Fatal("Error making IGDB covers request:", err)
	}

	var coverArt []CoverArt
	err = json.Unmarshal(body, &coverArt)
	if err != nil {
		log.Fatal("Error unmarshalling IGDB covers response:", err)
	}

	// Fetch accessibility information from PCGamingWiki
	// Closed captions
	closedCaptions := false
	colorBlind := false
	fullControllerSupport := false
	controllerRemapping := false
	if steamAvailability {
		resp, err := http.Get(fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Audio&fields=Audio.Closed_captions&join_on=Infobox_game._pageID=Audio._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID))
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading PC Gaming Wiki closed captions:", err)
		}
		var response PCGamingWikiResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal("Error unmarshalling PC Gaming Wiki closed captions:", err)
		}
		if len(response.CargoQuery) != 0 {
			closedCaptions, err = strconv.ParseBool(response.CargoQuery[0].Title.ClosedCaptions)
			if err != nil {
				log.Fatal("Error parsing closed captions:", err)
			}
		} else {
			closedCaptions = false
		}

		// Color blind mode
		resp, err = http.Get(fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Video&fields=Video.Color_blind&join_on=Infobox_game._pageID=Video._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID))
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading PC Gaming Wiki color blind mode:", err)
		}

		var response2 PCGamingWikiResponse
		err = json.Unmarshal(body, &response2)
		if err != nil {
			log.Fatal("Error unmarshalling PC Gaming Wiki color blind mode:", err)
		}
		if len(response2.CargoQuery) != 0 {
			colorBlind, err = strconv.ParseBool(response2.CargoQuery[0].Title.ColorBlind)
			if err != nil {
				log.Fatal("Error parsing color blind mode:", err)
			}
		} else {
			colorBlind = false
		}

		// Controller support and remapping
		resp, err = http.Get(fmt.Sprintf("https://www.pcgamingwiki.com/w/api.php?action=cargoquery&tables=Infobox_game,Input&fields=Input.Full_controller_support,Input.Controller_remapping,&join_on=Infobox_game._pageID=Input._PageID&where=Infobox_game.Steam_AppID%%20HOLDS%%20%%22%s%%22&format=json", steamID))
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading PC Gaming Wiki controller support and remapping:", err)
		}
		var response3 PCGamingWikiResponse
		err = json.Unmarshal(body, &response3)
		if err != nil {
			log.Fatal("Error unmarshalling PC Gaming Wiki controller support and remapping:", err)
		}
		if len(response3.CargoQuery) != 0 {
			fullControllerSupport, err = strconv.ParseBool(response3.CargoQuery[0].Title.FullControllerSupport)
			if err != nil {
				log.Fatal("Error parsing full controller support:", err)
			}
			controllerRemapping, err = strconv.ParseBool(response3.CargoQuery[0].Title.ControllerRemapping)
			if err != nil {
				log.Fatal("Error parsing controller remapping:", err)
			}
		} else {
			fullControllerSupport = false
			controllerRemapping = false
		}
	}

	//set fetched data to game
	game[0].CoverArt = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", coverArt[0].ImageID)
	game[0].ClosedCaptions = closedCaptions
	game[0].ColorBlind = colorBlind
	game[0].FullControllerSupport = fullControllerSupport
	game[0].ControllerRemapping = controllerRemapping
	game[0].SteamAvailability = steamAvailability

	// Marshal game data to JSON
	body, err = json.Marshal(game)
	if err != nil {
		log.Fatal("Error marshalling game data:", err)
	}

	return body, nil

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
