package handlers

import (
	"fmt"
	"log"
	"net/http"
	"playability/pkg/calc"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (env *Env) GetScoreHandler(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game")
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		log.Println("Error converting game ID to int: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	scores, err := env.DB.QueryAccessibilityScores(gameIDInt)
	if err != nil {
		log.Println("Error querying accessibility scores:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(scores) < 1 {
		http.Error(w, "Not enough reports", http.StatusNotAcceptable)
		return
	} else {
		accessibilityScore := calc.CalculateAccessibilityScore(scores)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%f", accessibilityScore)
	}
}
