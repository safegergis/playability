package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"playability/pkg/calc"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (env *Env) GetFeatureReportsHandler(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game")
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		log.Println("Error converting game ID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reports, err := env.DB.QueryFeatureReports(gameIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(reports) < 1 {
		http.Error(w, "Not enough reports", http.StatusNotAcceptable)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(calc.CalculateFeatureScore(reports))
	}
}
