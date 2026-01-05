package interceptor

import (
	"context"
	"event-collector/internal/monitoring"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicRecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			monitoring.GrpcPanicsTotal.WithLabelValues(info.FullMethod).Inc()
			err = status.Error(codes.Internal, "internal server error")
			resp = nil
		}
	}()

	return handler(ctx, req)
}
