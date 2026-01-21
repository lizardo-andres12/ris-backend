package otel

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// NewLoggerProvider creates a new provider based on resource and exporter configuration
func NewLoggerProvider(resource *resource.Resource, exporter *otlploggrpc.Exporter) *log.LoggerProvider {
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewSimpleProcessor(exporter)),
		log.WithResource(resource),
	)
	global.SetLoggerProvider(lp)
	return lp
}

