package service

import (
	"currency-notifier/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bojanz/currency"
	"log"
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

type CurrencyService struct {
	// Cache variables
	cacheRate  float64
	cacheTime  time.Time
	cacheMutex sync.Mutex
	// Dependencies
	repo *repository.ExchangeRateRepository
}

func NewCurrencyService(repo *repository.ExchangeRateRepository) *CurrencyService {
	return &CurrencyService{
		repo: repo,
	}
}

func (s *CurrencyService) Init() error {
	// TODO - implement calling API. If API call fails - init cache with latest rate from the database
	return s.ReloadRate()
}

func (s *CurrencyService) ReloadRate() error {
	_, err := s.reloadRate()

	return err
}

func (s *CurrencyService) reloadRate() (float64, error) {
	rate, err := fetchRateFromAPI()
	if err != nil {
		return 0, err
	}
	log.Printf("USD to UAH rate: %f", rate)

	s.updateCache(rate)
	err = s.repo.SaveRate(rate)
	if err != nil {
		return 0, err
	}

	log.Printf("USD to UAH rate saved to the database")

	return rate, nil
}

func (s *CurrencyService) GetUSDtoUAHRate() (float64, error) {
	if rate, ok := s.getCachedRate(); ok {
		return rate, nil
	}

	return s.reloadRate()
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

func (s *CurrencyService) getCachedRate() (float64, bool) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if time.Since(s.cacheTime) < time.Hour {
		return s.cacheRate, true
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

func (s *CurrencyService) updateCache(rate float64) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cacheRate = rate
	s.cacheTime = time.Now()
}
