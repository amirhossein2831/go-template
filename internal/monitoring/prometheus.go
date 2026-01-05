package monitoring

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	configs "event-collector/internal/config"
	"log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

func NewMetricsServer(cfg *configs.Config) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.PrometheusConfig.Host, cfg.PrometheusConfig.Port)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
}

func RunMetricsServer(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Prometheus metrics endpoint started: %v, path: /metrics", srv.Addr)

			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Println("Failed to run prometheus server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Stopping Prometheus metrics server, address: %v", srv.Addr)
			return srv.Shutdown(ctx)
		},
	})
}
