package jobs

import (
	"currency-notifier/internal/service"
	"go.uber.org/zap"
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
	zap.S().Debugf("Reloading exchange rate")
	err := j.currencyService.ReloadRate()
	if err != nil {
		zap.S().Errorf("Error fetching exchange rate: %v", err)
		return
	}
}
