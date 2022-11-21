package provider

import (
	"context"
	"fmt"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/hsmtkk/bookish-pancake/constant"
	"github.com/hsmtkk/bookish-pancake/utilgcp"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Provider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	projectID, err := utilgcp.ProjectID(ctx)
	if err != nil {
		return nil, err
	}
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return nil, fmt.Errorf("trace.New failed; %w", err)
	}
	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(constant.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource.New failed; %w", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, err
}
