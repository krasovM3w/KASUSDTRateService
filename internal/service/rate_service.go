package service

import (
	"context"
	"finalyTask/internal/repository/postgres"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type RateServicer interface {
	GetCurrentRate(ctx context.Context, base, target string) (*postgres.Rate, error)
	GetRate(ctx context.Context, base, target string) (*postgres.Rate, error)
}

type RateService struct {
	repo   postgres.RateRepository
	client PriceProvider
	logger *zap.Logger
}

func NewRateService(repo postgres.RateRepository, client PriceProvider, logger *zap.Logger) *RateService {
	return &RateService{repo: repo, client: client, logger: logger}
}

func (r *RateService) GetRate(ctx context.Context, base, target string) (*postgres.Rate, error) {
	price, err := r.client.GetUSDTPrice()
	if err != nil {
		r.logger.Error("Failed to get USDT price", zap.String("base", base), zap.String("target", target))
		return nil, err
	}

	rate := &postgres.Rate{
		BaseCurrency:   base,
		TargetCurrency: target,
		Rate:           price,
		Timestamp:      time.Now().UTC(),
	}
	if err := r.repo.SaveRate(ctx, base, target, rate.Rate); err != nil {
		r.logger.Error("Failed to save rate", zap.String("base", base), zap.String("target", target))
	}

	return rate, nil

}

func (s *RateService) GetCurrentRate(ctx context.Context, base, target string) (*postgres.Rate, error) {
	symbol := base + target
	price, err := s.client.GetPrice(symbol)
	if err != nil {
		s.logger.Error("Failed to get price from MEXC",
			zap.String("symbol", symbol),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	rate := &postgres.Rate{
		BaseCurrency:   base,
		TargetCurrency: target,
		Rate:           price,
		Timestamp:      time.Now().UTC(),
	}

	if err := s.repo.SaveRate(ctx, base, target, price); err != nil {
		s.logger.Error("Failed to save rate",
			zap.String("base", base),
			zap.String("target", target),
			zap.Float64("rate", price),
			zap.Error(err))
		return nil, fmt.Errorf("failed to save rate: %w", err)
	}

	return rate, nil
}
