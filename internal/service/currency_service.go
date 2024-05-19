package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bojanz/currency"
	"net/http"
	"strconv"
	"sync"
	"time"
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

// Cache variables
var cacheRate float64
var cacheTime time.Time
var cacheMutex sync.Mutex

func GetUSDtoUAHRate() (float64, error) {
	if rate, ok := getCachedRate(); ok {
		return rate, nil
	}

	rate, err := fetchRateFromAPI()
	if err != nil {
		return 0, err
	}

	updateCache(rate)
	return rate, nil
}

func getCurrencyCode(currencyCode string) int {
	code, ok := currency.GetNumericCode(currencyCode)
	if !ok {
		// TODO log here
	}
	numericCode, err := strconv.Atoi(code)
	if err != nil {
		// TODO log here
	}

	return numericCode
}

func getCachedRate() (float64, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(cacheTime) < time.Hour {
		return cacheRate, true
	}
	return 0, false
}

func fetchRateFromAPI() (float64, error) {
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

func updateCache(rate float64) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cacheRate = rate
	cacheTime = time.Now()
}
