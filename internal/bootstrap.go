package internal

import (
	configs "event-collector/internal/config"
	"event-collector/internal/database"
	"event-collector/internal/pkg/logger"
	"event-collector/internal/services"
	"event-collector/internal/transport/grpc"
	handlers2 "event-collector/internal/transport/grpc/handlers"
	"event-collector/internal/transport/http"
	"event-collector/internal/transport/http/handlers"
	"event-collector/internal/transport/http/route"
	"event-collector/pkg/env"
	"strings"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
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
		fx.Provide(
			configs.NewConfig,
			logger.NewZapLogger,
			database.NewMongo,
			services.NewGreetingService,
			handlers.NewGreetingHandler,
			handlers2.NewGreetingHandler,
			http.NewHTTPServer,
			grpc.NewGRPCServer,
		),

		fx.Invoke(
			database.RunMigration,
			func(*fiber.App) {},
			route.RegisterRoutes,
			func(server *grpc2.Server) {},
			grpc.RegisterServices,
		),
	)

	app.Run()
}
