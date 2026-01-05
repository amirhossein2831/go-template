package interceptor

import (
	"context"
	"event-collector/internal/monitoring"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func GrpcResponseTimeMetricInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	timer := prometheus.NewTimer(monitoring.GrpcResponseTime.WithLabelValues(info.FullMethod))
	defer timer.ObserveDuration()

	return handler(ctx, req)
}
