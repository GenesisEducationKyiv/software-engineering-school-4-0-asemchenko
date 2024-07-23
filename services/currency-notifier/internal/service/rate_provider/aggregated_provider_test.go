package rate_provider

import (
	"errors"
	"testing"
)

type mockRateProvider struct {
	rate float64
	err  error
}

func (m *mockRateProvider) FetchRateFromAPI() (float64, error) {
	return m.rate, m.err
}

func (m *mockRateProvider) GetName() string {
	return "mock-rate-provider"
}

func TestAggregatedRateProvider_FetchRateFromAPI(t *testing.T) {
	tests := []struct {
		name      string
		providers []RateProvider
		want      float64
		wantErr   bool
	}{
		{
			name: "Test with multiple providers",
			providers: []RateProvider{
				&mockRateProvider{rate: 28.25, err: nil},
				&mockRateProvider{rate: 28.75, err: nil},
			},
			want:    28.25,
			wantErr: false,
		},
		{
			name: "Test with one provider returning an error",
			providers: []RateProvider{
				&mockRateProvider{rate: 0, err: errors.New("error fetching rate")},
				&mockRateProvider{rate: 28.25, err: nil},
			},
			want:    28.25,
			wantErr: false,
		},
		{
			name: "Test with all providers returning an error",
			providers: []RateProvider{
				&mockRateProvider{rate: 0, err: errors.New("error fetching rate")},
				&mockRateProvider{rate: 0, err: errors.New("error fetching rate")},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregatedProvider := NewAggregatedRateProvider(tt.providers...)
			got, err := aggregatedProvider.FetchRateFromAPI()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchRateFromAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchRateFromAPI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
