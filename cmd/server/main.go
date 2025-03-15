package main

import (
	"context"
	"database/sql"
	grpc2 "finalyTask/internal/api/grpc"
	"finalyTask/internal/config"
	logger2 "finalyTask/internal/logger"
	"finalyTask/internal/repository/postgres"
	"finalyTask/internal/service"
	"finalyTask/internal/telemetry"
	"finalyTask/proto/currensy/github.com/m3w/usdt-rate/proto/currensy"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	logger, err := logger2.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	//tracing
	tp, err := telemetry.InitTracing(cfg)
	if err != nil {
		logger.Fatal("Error initializing tracing", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatal("Error shutting down tracing", zap.Error(err))
		}
	}()

	//metrics
	telemetry.InitMetrics()

	//db
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		logger.Fatal("Error opening database connection", zap.Error(err))
	}
	//db retry
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", cfg.DBURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		logger.Info("Retrying database connection...", zap.Int("attempt", i+1))
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Fatal("Failed to connect to database after retries", zap.Error(err))
	}
	defer db.Close()

	//ping db
	if err := db.Ping(); err != nil {
		logger.Fatal("Error pinging database", zap.Error(err))
	}

	//init
	rateRepo := postgres.NewPostgresRepo(db)
	mexcClient := service.NewMexcClient(cfg)
	rateService := service.NewRateService(rateRepo, mexcClient, logger)
	handler := grpc2.NewHandler(db, rateService, logger)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(telemetry.MetricsUnaryInterceptor))

	currensy.RegisterCurrencyServiceServer(grpcServer, handler)

	//reflection API
	reflection.Register(grpcServer)

	//grpc start
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("Starting gRPC server", zap.String("address", lis.Addr().String()))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.MetricsPort),
		Handler: promhttp.Handler(),
	}

	go func() {
		logger.Info("Starting metrics server", zap.String("address", metricsServer.Addr))
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Metrics server failed", zap.Error(err))
		}
	}()

	fmt.Println("started...")

	//graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//stop metric
	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Error("Metrics server shutdown error", zap.Error(err))
	}

	//stop grpc
	grpcServer.GracefulStop()

	logger.Info("Servers stopped gracefully")
}
