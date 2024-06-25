package nbu

import (
	"testing"
)

func Test_findNbuUsdRate(t *testing.T) {
	type args struct {
		rates []nbuRateExchange
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Test with USD rate",
			args: args{
				rates: []nbuRateExchange{
					{
						CurrencyCode: "USD",
						Rate:         28.25,
					},
				},
			},
			want:    28.25,
			wantErr: false,
		},
		{
			name: "Test without USD rate",
			args: args{
				rates: []nbuRateExchange{
					{
						CurrencyCode: "EUR",
						Rate:         0.85,
					},
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findNbuUsdRate(tt.args.rates)
			if (err != nil) != tt.wantErr {
				t.Errorf("findNbuUsdRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findNbuUsdRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
