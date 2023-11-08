package telemetry

import (
	"context"

	"github.com/dose-na-nuvem/toll-station/config"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"
)

const (
	serviceName    = "pedagio"
	serviceVersion = "0.1.0"

	meterName = "pedagio"
)

type Telemetry struct {
	MeterProvider  metric.MeterProvider
	TrafficCounter metric.Int64Counter
	GateHistogram  metric.Int64Histogram
}

var (
	mp sdk.MeterProvider
	// exp *prometheus.exporter
)

func New(ctx context.Context) (*Telemetry, error) {

	res, err := newResource()
	if err != nil {
		return nil, err
	}

	err = setupMetrics(ctx, res)
	if err != nil {
		return nil, err
	}

	return &Telemetry{
		TrafficCounter: trafficCounter,
		GateHistogram:  gateLatency,
	}, nil
}

func (t Telemetry) Start(cfg *config.Cfg) {

	pullBasedEndpoint := cfg.Telemetry.Metrics.Endpoint
	cfg.Logger.Info("servindo métricas em ", zap.String("endpoint", pullBasedEndpoint))

	go func() {
		err := serveMetrics(pullBasedEndpoint, metricsHandler)
		if err != nil {
			slog.Error("falha ao servir métricas", err)
		}
	}()
}

func (t Telemetry) Shutdown(ctx context.Context) {
	// fecha o provedor de métricas
	if err := mp.Shutdown(ctx); err != nil {
		slog.Error("provedor de métricas finalizado com falha", err)
	}
}
