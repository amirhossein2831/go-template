package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// TODO: use your service name
var GrpcErrorsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "service_name_grpc_errors_total",
		Help: "Total number of gRPC handler errors (non-OK responses)",
	},
	[]string{"method", "code"},
)

var GrpcSuccessTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "message_go_grpc_success_total",
		Help: "Total number of successful gRPC handler calls (OK responses)",
	},
	[]string{"method"},
)

// TODO: use your service name
var GrpcPanicsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "service_name_grpc_panics_total",
		Help: "Total number of recovered panics across all gRPC handlers",
	},
	[]string{"method"},
)

// TODO: use your service name
var GrpcResponseTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "service_name_grpc_response_time",
		Buckets: []float64{0.0001, 0.0005, 0.001, 0.002, 0.005, 0.01, 0.015, 0.02, 0.03, 0.04, 0.05, 0.075, 0.1, 0.2, 0.5, 0.75, 1, 1.5, 2, 2.5, 5},
	},
	[]string{"method"},
)
