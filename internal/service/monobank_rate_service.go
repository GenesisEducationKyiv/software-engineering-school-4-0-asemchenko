package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bojanz/currency"
	"io"
	"log"
	"net/http"
	"strconv"
)

// ISO 4217 currency codes
var usdCode = getCurrencyCode("USD")
var uahCode = getCurrencyCode("UAH")

func fetchRateFromAPI() (float64, error) {
	resp, err := http.Get(monobankAPI)
	if err != nil {
		return 0, err
	}
	defer closeBody(resp.Body)

	rates, err := extractRates(resp)
	if err != nil {
		return 0, err
	}

	return findUahRate(rates)
}

func extractRates(resp *http.Response) ([]CurrencyRate, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rates []CurrencyRate
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, err
	}

	return rates, nil
}

func closeBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func findUahRate(rates []CurrencyRate) (float64, error) {
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
