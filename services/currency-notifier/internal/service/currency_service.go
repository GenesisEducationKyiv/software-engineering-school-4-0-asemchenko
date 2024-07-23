package service

import (
	"currency-notifier/internal/service/rate_provider"
	"currency-notifier/internal/util"
	"go.uber.org/zap"
	"time"
)

type UpdateRateSubscriber interface {
	PublishRateUpdated(usdUahRate float64) error
}

type CurrencyService struct {
	repo         ExchangeRateRepo
	rateProvider rate_provider.RateProvider
	latestRate   *util.InMemoryCache[float64]
	subscribers  []UpdateRateSubscriber
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

func (s *CurrencyService) AddSubscriber(subscriber UpdateRateSubscriber) {
	s.subscribers = append(s.subscribers, subscriber)
}

func (s *CurrencyService) reloadRate() (float64, error) {
	rate, err := s.rateProvider.FetchRateFromAPI()
	if err != nil {
		return 0, err
	}
	saveRate, err := s.saveRate(rate)
	if err != nil {
		return 0, err
	}

	for _, subscriber := range s.subscribers {
		err = subscriber.PublishRateUpdated(saveRate)
		if err != nil {
			zap.S().Warnf("Failed to notify subscriber: %v", err)
		}
	}

	return saveRate, err
}

func (s *CurrencyService) saveRate(rate float64) (float64, error) {
	zap.S().Debugf("USD to UAH rate: %f", rate)

	s.latestRate.Set(rate)

	err := s.repo.SaveRate(rate)
	if err != nil {
		return 0, err
	}

	zap.S().Debugf("USD to UAH rate saved to the database")
	return rate, nil
}

func (s *CurrencyService) GetUSDtoUAHRate() (float64, error) {
	if rate, ok := s.latestRate.Get(); ok {
		return rate, nil
	}
	return s.reloadRate()
}
