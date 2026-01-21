package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// NewMetricProvider creates a new provider based on resource and exporter configuration
func NewMetricProvider(resource *resource.Resource, exporter *otlpmetricgrpc.Exporter) *metric.MeterProvider {
	mp := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(mp)
	return mp
}

