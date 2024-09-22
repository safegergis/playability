package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"playability/types"

	"github.com/go-chi/jwtauth/v5"
)

func (env *Env) PostReportHandler(w http.ResponseWriter, r *http.Request) {
	var reportBody types.Report
	if err := json.NewDecoder(r.Body).Decode(&reportBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("ReportBody: ", reportBody)

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		fmt.Println("Error getting claims: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims["sub"].(string)

	report := types.ReportRow{
		GameID:                reportBody.GameID,
		UserID:                userID,
		ClosedCaptions:        reportBody.ClosedCaptions,
		ColorBlind:            reportBody.ColorBlind,
		FullControllerSupport: reportBody.FullControllerSupport,
		ControllerRemapping:   reportBody.ControllerRemapping,
		Report:                reportBody.Report,
		Score:                 reportBody.Score,
	}

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

	w.WriteHeader(http.StatusCreated)

}

func (env *Env) GetReportsHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := env.DB.GetReports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}
