package telemetry

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ressourceName string) (shutdown func()) {
	ctx := context.Background()

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(ressourceName)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting otel-collector")
	}
	otel.SetTracerProvider(tracerProvider)

	shutdown = func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Error().Msg("Error cleaning container")
		}
	}

	return shutdown
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ressourceName string) (*trace.TracerProvider, error) {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("127.0.0.1:4317"),
			otlptracegrpc.WithRetry(
				otlptracegrpc.RetryConfig{
					Enabled: false,
				},
			),
		),
	)
	if err != nil {
		log.Debug().Err(err).Msgf("Could not set exporter: %v", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(ressourceName),
			semconv.TelemetrySDKLanguageGo,
		),
	)
	if err != nil {
		log.Debug().Err(err).Msgf("Could not set resources: %v", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(resources),
	)
	return traceProvider, nil
}
