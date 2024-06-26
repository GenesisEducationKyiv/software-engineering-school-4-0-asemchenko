// Code generated by MockGen. DO NOT EDIT.
// Source: ./rate_provider/rate_provider.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRateProvider is a mock of RateProvider interface.
type MockRateProvider struct {
	ctrl     *gomock.Controller
	recorder *MockRateProviderMockRecorder
}

// MockRateProviderMockRecorder is the mock recorder for MockRateProvider.
type MockRateProviderMockRecorder struct {
	mock *MockRateProvider
}

// NewMockRateProvider creates a new mock instance.
func NewMockRateProvider(ctrl *gomock.Controller) *MockRateProvider {
	mock := &MockRateProvider{ctrl: ctrl}
	mock.recorder = &MockRateProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateProvider) EXPECT() *MockRateProviderMockRecorder {
	return m.recorder
}

// FetchRateFromAPI mocks base method.
func (m *MockRateProvider) FetchRateFromAPI() (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchRateFromAPI")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchRateFromAPI indicates an expected call of FetchRateFromAPI.
func (mr *MockRateProviderMockRecorder) FetchRateFromAPI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchRateFromAPI", reflect.TypeOf((*MockRateProvider)(nil).FetchRateFromAPI))
}

// GetName mocks base method.
func (m *MockRateProvider) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockRateProviderMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockRateProvider)(nil).GetName))
}
