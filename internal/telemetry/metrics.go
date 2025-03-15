package telemetry

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"time"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_response_time_seconds",
			Help:    "Response time of gRPC methods",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(requestCounter, responseTime)
}

func MetricsUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	status := "success"
	if err != nil {
		status = "error"
	}

	requestCounter.WithLabelValues(info.FullMethod, status).Inc()
	responseTime.WithLabelValues(info.FullMethod).Observe(duration.Seconds())

	return resp, err
}
