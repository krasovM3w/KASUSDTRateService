package grpc

import (
	"context"
	"testing"
	"time"

	"finalyTask/internal/repository/postgres"
	"finalyTask/proto/currensy/github.com/m3w/usdt-rate/proto/currensy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetRate(ctx context.Context, base, target string) (*postgres.Rate, error) {
	args := m.Called(ctx, base, target)
	return args.Get(0).(*postgres.Rate), args.Error(1)
}

func (m *MockRateService) GetCurrentRate(ctx context.Context, base, target string) (*postgres.Rate, error) {
	args := m.Called(ctx, base, target)
	return args.Get(0).(*postgres.Rate), args.Error(1)
}

func TestHandler_GetRate(t *testing.T) {
	// Arrange
	mockService := new(MockRateService)
	logger, _ := zap.NewDevelopment()

	handler := NewHandler(nil, mockService, logger)

	expectedRate := &postgres.Rate{
		BaseCurrency:   "USDT",
		TargetCurrency: "USD",
		Rate:           1.0,
		Timestamp:      time.Now(),
	}

	mockService.On("GetCurrentRate", mock.Anything, "USDT", "USD").
		Return(expectedRate, nil)

	// Act
	resp, err := handler.GetRate(context.Background(), &currensy.GetRateRequest{
		BaseCurrency:   "USDT",
		TargetCurrency: "USD",
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRate.Rate, resp.Rate)
	assert.Equal(t, expectedRate.Timestamp.Format(time.RFC3339), resp.Timestamp)
	mockService.AssertExpectations(t)
}
