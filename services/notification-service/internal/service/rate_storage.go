package service

import (
	"errors"
	"go.uber.org/zap"
)

var (
	ErrRateNotFound = errors.New("rate not found")
)

type RateStorage interface {
	SaveRate(rateCode string, rate float64)
	GetLatestRate() (float64, error)
}

type storage struct {
	rates map[string]float64
}

func NewRateStorage() RateStorage {
	return &storage{
		rates: make(map[string]float64),
	}
}

func (s *storage) SaveRate(rateCode string, rate float64) {
	zap.S().Debugf("Saving rate %s: %f", rateCode, rate)
	s.rates[rateCode] = rate
}

func (s *storage) GetLatestRate() (float64, error) {
	rate, ok := s.rates["usd_uah"]
	if !ok {
		return 0, ErrRateNotFound
	}

	return rate, nil
}
