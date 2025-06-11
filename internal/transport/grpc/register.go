package grpc

import (
	"event-collector/internal/config"
	"event-collector/internal/transport/grpc/handler"
	providev1 "event-collector/pkg/grpc/provide/v1"
	"event-collector/pkg/parse"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServiceParams defines all dependencies needed for registration using a fx.In struct.
// fx will automatically populate the fields of this struct.
type ServiceParams struct {
	fx.In
	Server          *grpc.Server
	GreetingHandler *handler.GreetingHandler
	// As you add more handlers, you simply add new fields here.
	// e.g., OtherHandler *other.Handler
}

// RegisterServices is a fx.Invoke function that registers all gRPC handlers.
// Its signature is now clean and will not change as you add more services.
func RegisterServices(p ServiceParams, cfg *config.Config) {
	providev1.RegisterGreetingServiceServer(p.Server, p.GreetingHandler)
	// e.g., otherv1.RegisterOtherServiceServer(p.Server, p.OtherHandler)

	// Enable gRPC server reflection. allow query to server and discover available services
	if needReflection, err := parse.ToPrimary[bool](cfg.GetEnv(cfg.Server.GRPC.Reflection)); err == nil && needReflection {
		reflection.Register(p.Server)
	}
}
