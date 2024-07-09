package rate_provider

import (
	"fmt"
	"log"
)

type aggregatedRateProvider struct {
	providers []RateProvider
}

func NewAggregatedRateProvider(providers ...RateProvider) RateProvider {
	return &aggregatedRateProvider{
		providers: providers,
	}
}

func (s *aggregatedRateProvider) FetchRateFromAPI() (float64, error) {
	// iterate over all providers and return the first successful result
	for _, provider := range s.providers {
		rate, err := getRateFromProvider(provider)
		if err == nil {
			return rate, nil
		}
	}

	return 0, fmt.Errorf("all providers failed to fetch rate")
}

func (s *aggregatedRateProvider) GetName() string {
	return "aggregated-rate-provider"
}

func getRateFromProvider(provider RateProvider) (float64, error) {
	rate, err := provider.FetchRateFromAPI()
	if err != nil {
		log.Printf("Failed to fetch rate from %s: %s\n", provider.GetName(), err)
	} else {
		log.Printf("Fetched rate from %s: %f\n", provider.GetName(), rate)
	}

	return rate, err
}
