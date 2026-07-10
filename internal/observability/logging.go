package observability

import (
	"context"
	"fmt"
	"log/slog"
	"order-api/internal/config"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type Shutdown func(context.Context) error

func InitLogging(ctx context.Context, cfg config.TelemetryConfig) (*slog.Logger, Shutdown, error) {
	exporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpoint(cfg.OTLPEndpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuat OTLP log exporter: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuat resource: %w", err)
	}

	processor := sdklog.NewSimpleProcessor(exporter)

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(res),
	)

	global.SetLoggerProvider(provider)

	logger := otelslog.NewLogger(cfg.ServiceName)

	slog.SetDefault(logger)

	return logger, provider.Shutdown, nil

}
