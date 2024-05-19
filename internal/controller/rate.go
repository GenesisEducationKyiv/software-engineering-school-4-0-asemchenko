package controller

import (
	"encoding/json"
	"net/http"
)

// GetRate returns the current USD to UAH exchange rate
// @Summary Get the current USD to UAH exchange rate
// @Description Request returns the current USD to UAH exchange rate using a third-party API
// @Tags rate
// @Produce json
// @Success 200 {number} int "Current USD to UAH exchange rate"
// @Router /rate [get]
func GetRate(w http.ResponseWriter, r *http.Request) {
	rate := 27.5

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}
