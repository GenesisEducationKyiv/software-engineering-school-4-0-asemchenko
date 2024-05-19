package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bojanz/currency"
	"net/http"
	"strconv"
)

const monobankAPI = "https://api.monobank.ua/bank/currency"

type CurrencyRate struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	Date          int64   `json:"date"`
	RateSell      float64 `json:"rateSell,omitempty"`
	RateBuy       float64 `json:"rateBuy,omitempty"`
	RateCross     float64 `json:"rateCross,omitempty"`
}

// ISO 4217 currency codes
var usdCode = getCurrencyCode("USD")
var uahCode = getCurrencyCode("UAH")

func GetUSDtoUAHRate() (float64, error) {
	resp, err := http.Get(monobankAPI)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rates []CurrencyRate
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return 0, err
	}

	for _, rate := range rates {
		if rate.CurrencyCodeA == usdCode && rate.CurrencyCodeB == uahCode {
			return rate.RateSell, nil
		}
	}

	return 0, errors.New("USD to UAH rate not found")
}

func getCurrencyCode(currencyCode string) int {
	code, ok := currency.GetNumericCode(currencyCode)
	if ok != true {
		// TODO log here
	}
	numericCode, err := strconv.Atoi(code)
	if err != nil {
		// TODO log here
	}

	return numericCode
}
