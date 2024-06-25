package service

import (
	"currency-notifier/internal/service/rate_provider"
	"currency-notifier/internal/util"
	"log"
	"time"
)

type CurrencyService struct {
	repo         ExchangeRateRepo
	rateProvider rate_provider.RateProvider
	latestRate   *util.InMemoryCache[float64]
}

type ExchangeRateRepo interface {
	SaveRate(rate float64) error
	GetLatestRate() (float64, error)
}

func NewCurrencyService(repo ExchangeRateRepo, rateProvider rate_provider.RateProvider) *CurrencyService {
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
