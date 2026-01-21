package otel

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

const (
	environmentResourceKey = "environment"
)

// NewResource creates the base resource object shared by all OpenTelemetry providers
func NewResource() (*resource.Resource, error) {
	version := os.Getenv("VERSION")
	environment := os.Getenv("ENVIRONMENT") // production, staging, etc.
	serviceName := os.Getenv("SERVICE_NAME")

	return resource.New(
		context.Background(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithContainer(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
			attribute.String(environmentResourceKey, environment),
		),
	)
}

