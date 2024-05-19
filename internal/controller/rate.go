package controller

import (
	"currency-notifier/internal/service"
	"encoding/json"
	"net/http"
)

// GetRate returns the current USD to UAH exchange rate
// @Summary Get the current USD to UAH exchange rate
// @Description Request returns the current USD to UAH exchange rate using Monobank API
// @Tags rate
// @Produce json
// @Success 200 {number} float64 "Current USD to UAH exchange rate"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rate [get]
func GetRate(w http.ResponseWriter, r *http.Request) {
	rate, err := service.GetUSDtoUAHRate()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}
