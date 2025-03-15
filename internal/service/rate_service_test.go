package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"testing"

	pg "finalyTask/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// PriceProvider
type MockPriceProvider struct {
	mock.Mock
}

func (m *MockPriceProvider) GetUSDTPrice() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPriceProvider) GetPrice(symbol string) (float64, error) {
	args := m.Called(symbol)
	return args.Get(0).(float64), args.Error(1)
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveRate(ctx context.Context, base, target string, rate float64) error {
	args := m.Called(ctx, base, target, rate)
	return args.Error(0)
}

func (m *MockRepo) GetLatestRates(ctx context.Context, base, target string, limit int) ([]pg.Rate, error) {
	args := m.Called(ctx, base, target, limit)
	return args.Get(0).([]pg.Rate), args.Error(1)
}

func TestRateService_GetCurrentRate_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepo)
	mockPriceProvider := new(MockPriceProvider)
	logger := zap.NewNop()

	service := NewRateService(mockRepo, mockPriceProvider, logger)

	mockPriceProvider.On("GetPrice", "KASUSDT").Return(1.0, nil)
	mockRepo.On("SaveRate", mock.Anything, "KAS", "USDT", 1.0).Return(nil)

	rate, err := service.GetCurrentRate(context.Background(), "KAS", "USDT")

	assert.NoError(t, err)
	assert.Equal(t, "KAS", rate.BaseCurrency)
	assert.Equal(t, "USDT", rate.TargetCurrency)
	assert.Equal(t, 1.0, rate.Rate)
	mockPriceProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestRateService_GetCurrentRate_APIError(t *testing.T) {
	mockRepo := new(MockRepo)
	mockPriceProvider := new(MockPriceProvider)
	logger := zap.NewNop()

	service := NewRateService(mockRepo, mockPriceProvider, logger)

	mockPriceProvider.On("GetPrice", "KASUSDT").Return(0.0, errors.New("API unavailable"))

	_, err := service.GetCurrentRate(context.Background(), "KAS", "USDT")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API unavailable")
	mockPriceProvider.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "SaveRate")
}

func TestRateService_GetUSDTPrice(t *testing.T) {
	mockRepo := new(MockRepo)
	mockPriceProvider := new(MockPriceProvider)
	logger := zap.NewNop()

	service := NewRateService(mockRepo, mockPriceProvider, logger)

	mockPriceProvider.On("GetPrice", "USDTUSD").Return(1.0, nil)
	mockRepo.On("SaveRate", mock.Anything, "USDT", "USD", 1.0).Return(nil)

	rate, err := service.GetCurrentRate(context.Background(), "USDT", "USD")
	assert.NoError(t, err)
	assert.Equal(t, 1.0, rate.Rate)

	mockPriceProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
