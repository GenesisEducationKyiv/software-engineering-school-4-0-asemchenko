package rate_provider

type RateProvider interface {
	FetchRateFromAPI() (float64, error)
}
