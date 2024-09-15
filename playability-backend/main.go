package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"bytes"
)

type Game struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Cover             int    `json:"cover"`
	Summary           string `json:"summary"`
	Platforms         []int  `json:"platforms"`
	InvolvedCompanies []int  `json:"involved_companies"`
	CoverArt          string `json:"cover_art"`
}
type CoverArt struct {
	ID      int    `json:"id"`
	ImageID string `json:"image_id"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/search", searchHandler).Methods("GET")
	router.HandleFunc("/games", gamesHandler).Methods("GET")

	http.ListenAndServe(":8080", router)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")

	postBody := fmt.Sprintf("fields id,game,name; search \"%s\"; limit 50;", searchTerm)

	body, err := makeIgdbRequest(postBody, igdbSecret, "search")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
	w.Write([]byte(body))
	w.Header().Set("Content-Type", "application/json")

}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	body, err := getGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func getGame(gameID string) ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	igdbSecret := os.Getenv("IGDB_ACCESS_TOKEN")

	postBody := fmt.Sprintf("fields name,cover,summary,platforms,involved_companies; where id = %s;", gameID)

	body, err := makeIgdbRequest(postBody, igdbSecret, "games")
	if err != nil {
		return nil, fmt.Errorf("error making IGDB games request: %w", err)
	}

	var game []Game
	err = json.Unmarshal(body, &game)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling IGDB games response: %w", err)
	}

	if len(game) == 0 {
		return nil, fmt.Errorf("no game found with ID %s", gameID)
	}

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

	game[0].CoverArt = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", coverArt[0].ImageID)

	body, err = json.Marshal(game)
	if err != nil {
		log.Fatal("Error marshalling game data:", err)
	}
	log.Println(string(body))
	return body, nil

}

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
