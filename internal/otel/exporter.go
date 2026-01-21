package otel

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

// NewTraceExporter creates and returns the base otlp trace exporter targeting OTel Collector
func NewTraceExporter() (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(context.Background())
}

// NewMetricExporter creates and returns the base otlp metric exporter tageting OTel collector
func NewMetricExporter() (*otlpmetricgrpc.Exporter, error) {
	return otlpmetricgrpc.New(context.Background())
}

// NewLogExporter creates and returns the base otlp log exporter targeting OTel collector
func NewLogExporter() (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(context.Background())
}

