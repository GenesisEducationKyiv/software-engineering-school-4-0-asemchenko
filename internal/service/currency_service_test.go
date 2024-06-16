package service

import (
	"currency-notifier/internal/mocks"
	"errors"
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

	assertRate(t, err, 8.0, rate)
}

func TestCurrencyService_ReloadRate_whenCachedValue_shouldReloadRateAndSaveIt(t *testing.T) {
	// prepare test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRateProvider := mocks.NewMockRateProvider(ctrl)
	mockExchangeRepo := mocks.NewMockExchangeRateRepo(ctrl)

	mockRateProvider.EXPECT().FetchRateFromAPI().Return(8.00, nil)

	currencyService := NewCurrencyService(mockExchangeRepo, mockRateProvider)
	// cache value 8.0 before the test
	mockExchangeRepo.EXPECT().SaveRate(8.0).Times(1)
	_, err := currencyService.GetUSDtoUAHRate()
	assertNoError(t, err)

	// perform test
	mockRateProvider.EXPECT().FetchRateFromAPI().Return(40.9, nil)
	rate, err := currencyService.GetUSDtoUAHRate()
	assertRate(t, err, 8.0, rate)

	mockExchangeRepo.EXPECT().SaveRate(40.9).Times(1)
	assertNoError(t, currencyService.ReloadRate())
	rate, err = currencyService.GetUSDtoUAHRate()
	assertRate(t, err, 40.9, rate)
}

func TestCurrencyService_ReloadRate_whenErrorFetchingRate_shouldReturnError(t *testing.T) {
	// prepare test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRateProvider := mocks.NewMockRateProvider(ctrl)
	mockExchangeRepo := mocks.NewMockExchangeRateRepo(ctrl)

	assertError := errors.New("error fetching rate")
	mockRateProvider.EXPECT().FetchRateFromAPI().Return(0.0, assertError)
	currencyService := NewCurrencyService(mockExchangeRepo, mockRateProvider)

	// perform test
	err := currencyService.ReloadRate()
	if !errors.Is(err, assertError) {
		t.Fatalf("Expected error %v, got %v", assertError, err)
	}

}

func assertRate(t *testing.T, err error, expected float64, actual float64) {
	// assert result
	assertNoError(t, err)
	if actual != expected {
		t.Fatalf("Expected actual %f, got %f", expected, actual)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
