package controller

import (
	"currency-notifier/internal/service"
	"currency-notifier/internal/util"
	"net/http"
)

type RateController struct {
	service *service.CurrencyService
}

func NewRateController(service *service.CurrencyService) *RateController {
	return &RateController{
		service: service,
	}
}

// GetRate returns the current USD to UAH exchange rate
// @Summary Get the current USD to UAH exchange rate
// @Description Request returns the current USD to UAH exchange rate using Monobank API
// @Tags rate
// @Produce json
// @Success 200 {number} float64 "Current USD to UAH exchange rate"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rate [get]
func (c *RateController) GetRate(w http.ResponseWriter, _ *http.Request) {
	rate, err := c.service.GetUSDtoUAHRate()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	util.RespondJSON(&w, http.StatusOK, rate)
}
