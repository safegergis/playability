package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"playability/pkg/ai"
	"playability/types"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

// PostReportHandler handles the creation of a new accessibility report
func (env *Env) PostReportHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request body
	var reportBody types.ReportRegister
	if err := json.NewDecoder(r.Body).Decode(&reportBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Extract user ID from JWT claims
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		fmt.Println("[PostReportHandler] Error getting claims: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if token == nil {
		fmt.Println("[PostReportHandler] Token is nil")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Printf("[PostReportHandler] Token: %+v, Claims: %+v\n", token, claims)

	// Extract user ID from claims (stored as string)
	userIDStr := claims["sub"].(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Println("[PostReportHandler] Error converting user ID: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Create a ReportRow struct with the report data
	report := types.ReportRow{
		GameID:                reportBody.GameID,
		UserID:                userID,
		Platform:              reportBody.Platform,
		ClosedCaptions:        reportBody.ClosedCaptions,
		ColorBlind:            reportBody.ColorBlind,
		FullControllerSupport: reportBody.FullControllerSupport,
		ControllerRemapping:   reportBody.ControllerRemapping,
		Report:                reportBody.Report,
		Score:                 reportBody.Score,
	}
	//TODO: Change this to a queue system to prevent failed requests from killing system

	moderationResponse, err := ai.Moderation(&report)
	if err != nil {
		log.Println("Error moderating report: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var moderation types.ModerationResponse
	err = json.Unmarshal(moderationResponse, &moderation)
	if err != nil {
		log.Println("Error unmarshalling moderation response: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if moderation.Violation {
		http.Error(w, "Report contains violation "+strings.Join(moderation.Categories, ", ")+": "+moderation.Explanation, http.StatusForbidden)
		return
	} else {
		// Insert the report into the database
		err = env.DB.InsertReport(&report)
		if err != nil {
			if err.Error() == "report already exists" {
				http.Error(w, "report already exists", http.StatusConflict)
				return
			} else {
				fmt.Println("Error inserting report: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// GetReportCardsHandler retrieves report cards for a specific game
func (env *Env) GetReportCardsHandler(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game")
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		log.Println("Error converting game ID to int: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reports, err := env.DB.QueryReportCards(gameIDInt)
	if err != nil {
		log.Println("Error getting report cards: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}
