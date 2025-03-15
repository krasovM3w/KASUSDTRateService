package grpc

import (
	"context"
	"database/sql"
	"finalyTask/internal/service"
	"finalyTask/proto/currensy/github.com/m3w/usdt-rate/proto/currensy"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"time"
)

type Handler struct {
	db          *sql.DB
	rateService service.RateServicer
	currensy.UnimplementedCurrencyServiceServer
	logger *zap.Logger
}

func NewHandler(db *sql.DB, rs service.RateServicer, logger *zap.Logger) *Handler {
	return &Handler{
		db:          db,
		rateService: rs,
		logger:      logger,
	}
}

func (h *Handler) GetRate(ctx context.Context, req *currensy.GetRateRequest) (*currensy.GetRateResponse, error) {
	ctx, span := otel.Tracer("currency-service").Start(ctx, "GetRate")
	defer span.End()
	rate, err := h.rateService.GetCurrentRate(ctx, req.BaseCurrency, req.TargetCurrency)
	if err != nil {
		h.logger.Error("failed to GetRate", zap.Error(err))
		return nil, fmt.Errorf("failed to GetRate: %w", err)
	}
	return &currensy.GetRateResponse{
		Rate:      rate.Rate,
		Timestamp: rate.Timestamp.Format(time.RFC3339),
	}, nil
}

func (h *Handler) HealthCheck(ctx context.Context, req *currensy.HealthCheckRequest) (*currensy.HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	h.logger.Info("HealthCheck called")

	if err := h.db.PingContext(ctx); err != nil {
		h.logger.Error("Database ping failed", zap.Error(err))
		return &currensy.HealthCheckResponse{
			Status: currensy.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	h.logger.Info("Database is healthy")
	return &currensy.HealthCheckResponse{
		Status: currensy.HealthCheckResponse_SERVING,
	}, nil
}
