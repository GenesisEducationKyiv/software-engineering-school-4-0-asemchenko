package jobs

import (
	"currency-notifier/internal/service"
	"log"
)

type UpdateRateJob struct {
	currencyService *service.CurrencyService
}

func NewUpdateRateJob(currencyService *service.CurrencyService) *UpdateRateJob {
	return &UpdateRateJob{
		currencyService: currencyService,
	}
}

func (j UpdateRateJob) UpdateExchangeRate() {
	log.Printf("Reloading exchange rate")
	err := j.currencyService.ReloadRate()
	if err != nil {
		log.Printf("Error fetching exchange rate: %v", err)
		return
	}
}
