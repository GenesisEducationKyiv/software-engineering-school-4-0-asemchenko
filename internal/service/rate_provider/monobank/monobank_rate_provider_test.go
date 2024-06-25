package monobank

import (
	"testing"
)

func Test_findUahRate(t *testing.T) {
	type args struct {
		rates []currencyRate
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Test with UAH rate",
			args: args{
				rates: []currencyRate{
					{
						CurrencyCodeA: usdCode,
						CurrencyCodeB: uahCode,
						RateSell:      28.25,
					},
				},
			},
			want:    28.25,
			wantErr: false,
		},
		{
			name: "Test without UAH rate",
			args: args{
				rates: []currencyRate{
					{
						CurrencyCodeA: usdCode,
						CurrencyCodeB: getCurrencyCode("EUR"),
						RateSell:      0.85,
					},
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findUahRate(tt.args.rates)
			if (err != nil) != tt.wantErr {
				t.Errorf("findUahRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findUahRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
