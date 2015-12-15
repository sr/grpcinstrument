package promeasurer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sr/operator/src/grpcinstrument"
	"go.pedge.io/proto/time"
)

type measurer struct {
	registry *registry
}

type registry struct {
	total    *prometheus.CounterVec
	errors   *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func New() grpcinstrument.Measurer {
	return &measurer{
		registry: &registry{
			total: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "grpc_calls_total",
				Help: "Number of gRPC calls received by the server being instrumented.",
			}, []string{"service", "method"}),
			errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "grpc_calls_errors",
				Help: "Number of gRPC calls that returned an error.",
			}, []string{"service", "method"}),
			duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "grpc_calls_durations",
				Help: "Duration of gRPC calls.",
			}, []string{"service", "method"}),
		},
	}
}

func (m *measurer) Init() error {
	if err := prometheus.Register(m.registry.total); err != nil {
		return err
	}
	if err := prometheus.Register(m.registry.errors); err != nil {
		return err
	}
	if err := prometheus.Register(m.registry.duration); err != nil {
		return err
	}
	return nil
}

func (m *measurer) Measure(call *grpcinstrument.Call) {
	labels := prometheus.Labels{"service": call.Service, "method": call.Method}
	m.registry.total.With(labels).Inc()
	m.registry.errors.With(labels).Inc()
	m.registry.duration.With(labels).Observe(
		float64(prototime.DurationFromProto(call.Duration).Nanoseconds()),
	)
}
