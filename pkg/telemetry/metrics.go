package telemetry

import (
	"context"
	"net/http"
	"runtime"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	metricsHandler http.Handler
	trafficCounter api.Int64Counter
	gateLatency    api.Int64Histogram
	registry       *prom.Registry
)

func newMeterProvider(res *resource.Resource, exp *prometheus.Exporter) *metric.MeterProvider {
	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exp),
	)
}

func newMetricsExporter() (*prometheus.Exporter, error) {
	registry = getPrometheusRegistry()

	opts := make([]prometheus.Option, 0)
	opts = append(opts, prometheus.WithRegisterer(registry))

	// Esse exporter de OpenTelemetry está embutido com um Reader OpenTelemetry
	// e também implementa o prometheus.Collector, permitindo ser usado tanto
	// como Leitor/Reader e Collector.
	return prometheus.New(opts...)
}

// Encapsula obtenção do registry do Prometheus sem a instrumentação de processo
func getPrometheusRegistry() *prom.Registry {
	return prom.NewRegistry()
}

func newTrafficCounter(meter api.Meter) (api.Int64Counter, error) {
	return meter.Int64Counter("traffic",
		api.WithDescription("Qtde de carros passando pela faixa sem-parar"),
		api.WithUnit("cars"),
	)
}

func newGateLatency(meter api.Meter) (api.Int64Histogram, error) {
	return meter.Int64Histogram("gate_opening_duration",
		api.WithDescription("tempo para decidir abertura da cancela"),
		api.WithUnit("ms"),
	)
}

func serveMetrics(endpoint string, handler http.Handler) error {

	mux := http.NewServeMux()
	mux.Handle("/metrics", handler)

	srv := &http.Server{
		Addr:              endpoint,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Registra um callback para fechar o server quando o processo terminar.
	runtime.SetFinalizer(srv, func(*http.Server) {
		srv.Close()
	})

	return srv.ListenAndServe()
}

func setupMetrics(ctx context.Context, res *resource.Resource) error {

	metricsExporter, err := newMetricsExporter()
	if err != nil {
		return err
	}

	mp := newMeterProvider(res, metricsExporter)
	otel.SetMeterProvider(mp)

	meter := otel.GetMeterProvider().Meter(meterName)
	trafficCounter, err = newTrafficCounter(meter)
	if err != nil {
		return err
	}

	gateLatency, err = newGateLatency(meter)
	if err != nil {
		return err
	}

	metricsHandler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return nil
}
