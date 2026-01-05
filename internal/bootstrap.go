package internal

import (
	configs "event-collector/internal/config"
	"event-collector/internal/database/mongo"
	"event-collector/internal/monitoring"
	"event-collector/internal/services"
	"event-collector/internal/transport/grpc"
	grpcHandler "event-collector/internal/transport/grpc/handlers"
	"event-collector/internal/transport/http"
	"event-collector/internal/transport/http/handlers"
	"event-collector/internal/transport/http/route"
	"event-collector/pkg/env"
	"event-collector/pkg/logger"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func Bootstrap() {
	serviceMode := env.GetEnv("SERVICE_MODE", "server")

	switch {
	case strings.EqualFold(serviceMode, string(configs.ServiceModeServer)):
		runServer()
	default:
		panic("unknown SERVICE_MODE: " + serviceMode)
	}
}

func runServer() {
	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		fx.Provide(
			configs.NewConfig,
			logger.NewZapLogger,
			monitoring.NewMetricsServer,
			mongo.NewMongo,
			services.NewGreetingService,
			handlers.NewGreetingHandler,
			grpcHandler.NewGreetingHandler,
			http.NewHTTPServer,
			grpc.NewGRPCServer,
		),

		fx.Invoke(
			monitoring.RunMetricsServer,
			mongo.RunMigration,
			route.RegisterRoutes,
			grpc.RegisterServices,
		),
	)

	app.Run()
}
