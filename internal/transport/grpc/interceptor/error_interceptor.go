package interceptor

import (
	"context"
	"event-collector/internal/monitoring"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GrpcStatusMetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		monitoring.GrpcErrorsTotal.WithLabelValues(info.FullMethod, st.Code().String()).Inc()
		return resp, err
	}

	monitoring.GrpcSuccessTotal.WithLabelValues(info.FullMethod).Inc()
	return resp, nil
}
