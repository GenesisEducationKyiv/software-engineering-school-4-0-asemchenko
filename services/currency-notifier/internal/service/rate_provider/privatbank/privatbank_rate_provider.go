package privatbank

import (
	"currency-notifier/internal/service/rate_provider"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type rateProvider struct {
	privatBankHostUrl string
}

func NewPrivatBankRateProvider(privatBankHostUrl string) rate_provider.RateProvider {
	return &rateProvider{
		privatBankHostUrl: privatBankHostUrl,
	}
}

type exchangeRate struct {
	Currency     string `json:"ccy"`
	BaseCurrency string `json:"base_ccy"`
	BuyRate      string `json:"buy"`
	SaleRate     string `json:"sale"`
}

func (p *rateProvider) FetchRateFromAPI() (float64, error) {
	resp, err := http.Get(p.privatBankHostUrl + "/p24api/pubinfo?exchange&coursid=11")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var rates []exchangeRate
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		return 0, err
	}
	log.Printf("privatbank rates: %v", rates)

	return findUahRate(rates)
}

func (p *rateProvider) GetName() string {
	return "privat-bank-rate-provider"
}

func findUahRate(rates []exchangeRate) (float64, error) {
	for _, rate := range rates {
		if rate.Currency == "USD" && rate.BaseCurrency == "UAH" {
			return parseRate(rate.BuyRate)
		}
	}
	return 0, errors.New("USD to UAH rate not found")
}

func parseRate(rate string) (float64, error) {
	return strconv.ParseFloat(rate, 64)
}
