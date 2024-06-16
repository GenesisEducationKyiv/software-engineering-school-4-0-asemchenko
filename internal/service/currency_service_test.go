package service

import (
	"currency-notifier/internal/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

//go:generate mockgen -source=./currency_service.go -destination=../mocks/currency_service.go -package=mocks

func TestCurrencyService_GetUSDtoUAHRate_whenNoCachedValue_shouldLoadAndSaveRate(t *testing.T) {
	// prepare test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRateProvider := mocks.NewMockRateProvider(ctrl)
	mockExchangeRepo := mocks.NewMockExchangeRateRepo(ctrl)

	mockRateProvider.EXPECT().FetchRateFromAPI().Return(8.00, nil)
	mockExchangeRepo.EXPECT().SaveRate(8.00).Times(1)

	currencyService := NewCurrencyService(mockExchangeRepo, mockRateProvider)

	// perform test
	rate, err := currencyService.GetUSDtoUAHRate()

	// assert result
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if rate != 8.0 {
		t.Fatalf("Expected rate 8.0, got %f", rate)
	}
}
