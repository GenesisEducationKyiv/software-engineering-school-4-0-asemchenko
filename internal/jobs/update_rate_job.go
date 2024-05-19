package jobs

import (
	"currency-notifier/internal/service"
	"log"
)

func UpdateExchangeRateJob(currencyService *service.CurrencyService) {
	log.Printf("Reloading exchange rate")
	err := currencyService.ReloadRate()
	if err != nil {
		log.Printf("Error fetching exchange rate: %v", err)
		return
	}
}
