package nbu

import (
	"currency-notifier/internal/service/rate_provider"
	"currency-notifier/internal/util"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type rateProvider struct {
	nbuHostUrl string
}

func NewNbuRateProvider(nbuHostUrl string) rate_provider.RateProvider {
	return &rateProvider{
		nbuHostUrl: nbuHostUrl,
	}
}

type nbuRateExchange struct {
	Text         string  `json:"txt"`
	CurrencyCode string  `json:"cc"`
	ExchangeDate string  `json:"exchangedate"`
	Rate         float64 `json:"rate"`
	R030         int     `json:"r030"`
}

func (p *rateProvider) FetchRateFromAPI() (float64, error) {
	resp, err := http.Get(p.nbuHostUrl + "/NBUStatService/v1/statdirectory/exchange?json")
	if err != nil {
		return 0, err
	}
	defer util.CloseBodyWithErrorHandling(resp.Body)

	var rates []nbuRateExchange
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		return 0, err
	}
	log.Printf("nbu rates: %v", rates)

	return findNbuUsdRate(rates)
}

func (p *rateProvider) GetName() string {
	return "nbu-rate-provider"
}

func findNbuUsdRate(rates []nbuRateExchange) (float64, error) {
	for _, rate := range rates {
		if rate.CurrencyCode == "USD" {
			return rate.Rate, nil
		}
	}
	return 0, errors.New("USD to UAH rate not found")
}
