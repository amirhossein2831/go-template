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

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	grpc2 "google.golang.org/grpc"
)

func Initialize() {
	serviceMode := env.GetEnv("SERVICE_MODE", "server")

	switch {
	case strings.EqualFold(serviceMode, string(env.ServiceModeServer)):
		InitializeServer()
	}
}

func InitializeServer() {
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
			func(*fiber.App) {},
			route.RegisterRoutes,
			func(server *grpc2.Server) {},
			grpc.RegisterServices,
		),
	)

	app.Run()
}
