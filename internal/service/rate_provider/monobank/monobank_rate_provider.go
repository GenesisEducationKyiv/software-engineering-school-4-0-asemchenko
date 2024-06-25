package monobank

import (
	"currency-notifier/internal/service/rate_provider"
	"currency-notifier/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bojanz/currency"
	"log"
	"net/http"
	"strconv"
)

// ISO 4217 currency codes
var usdCode = getCurrencyCode("USD")
var uahCode = getCurrencyCode("UAH")

type currencyRate struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	Date          int64   `json:"date"`
	RateSell      float64 `json:"rateSell,omitempty"`
	RateBuy       float64 `json:"rateBuy,omitempty"`
	RateCross     float64 `json:"rateCross,omitempty"`
}

type rateProvider struct {
	monobankHostUrl string
}

func NewMonobankRateProvider(monobankHostUrl string) rate_provider.RateProvider {
	return &rateProvider{
		monobankHostUrl: monobankHostUrl,
	}
}

func (p *rateProvider) FetchRateFromAPI() (float64, error) {
	resp, err := http.Get(p.monobankHostUrl + "/bank/currency")
	if err != nil {
		return 0, err
	}
	defer util.CloseBodyWithErrorHandling(resp.Body)

	rates, err := extractRates(resp)
	if err != nil {
		return 0, err
	}

	return findUahRate(rates)
}

func (p *rateProvider) GetName() string {
	return "monobank-rate-provider"
}

func extractRates(resp *http.Response) ([]currencyRate, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rates []currencyRate
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, err
	}
	log.Printf("monobank rates: %v", rates)

	return rates, nil
}

func findUahRate(rates []currencyRate) (float64, error) {
	for _, rate := range rates {
		if rate.CurrencyCodeA == usdCode && rate.CurrencyCodeB == uahCode {
			return rate.RateSell, nil
		}
	}

	return 0, errors.New("USD to UAH rate not found")
}

func getCurrencyCode(currencyCode string) int {
	code, ok := currency.GetNumericCode(currencyCode)
	if !ok {
		log.Fatal("Currency code not found")
	}

	numericCode, err := strconv.Atoi(code)
	if err != nil {
		log.Fatal("Currency code is not a number")
	}

	return numericCode
}
