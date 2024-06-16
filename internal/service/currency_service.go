package service

import (
	"currency-notifier/internal/util"
	"log"
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

type CurrencyService struct {
	repo         ExchangeRateRepo
	rateProvider RateProvider
	latestRate   *util.InMemoryCache[float64]
}

type ExchangeRateRepo interface {
	SaveRate(rate float64) error
	GetLatestRate() (float64, error)
}

type RateProvider interface {
	FetchRateFromAPI() (float64, error)
}

func NewCurrencyService(repo ExchangeRateRepo, rateProvider RateProvider) *CurrencyService {
	return &CurrencyService{
		repo:         repo,
		rateProvider: rateProvider,
		latestRate:   util.NewInMemoryCache[float64](time.Hour),
	}
}
func (s *CurrencyService) Init() error {
	return s.ReloadRate()
}

func (s *CurrencyService) ReloadRate() error {
	_, err := s.reloadRate()
	return err
}

func (s *CurrencyService) reloadRate() (float64, error) {
	rate, err := s.rateProvider.FetchRateFromAPI()
	if err != nil {
		return 0, err
	}
	return s.saveRate(rate)
}

func (s *CurrencyService) saveRate(rate float64) (float64, error) {
	log.Printf("USD to UAH rate: %f", rate)

	s.latestRate.Set(rate)

	err := s.repo.SaveRate(rate)
	if err != nil {
		return 0, err
	}

	log.Printf("USD to UAH rate saved to the database")
	return rate, nil
}

func (s *CurrencyService) GetUSDtoUAHRate() (float64, error) {
	if rate, ok := s.latestRate.Get(); ok {
		return rate, nil
	}
	return s.reloadRate()
}
