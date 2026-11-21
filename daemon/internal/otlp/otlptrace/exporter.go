package otlptrace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	tracev1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type (
	Client interface {
		UploadTraces(ctx context.Context, protoSpans []*tracev1.ResourceSpans) error
	}
	Exporter = otlptrace.Exporter
)

func New(ctx context.Context, client Client) (*Exporter, error) {
	return otlptrace.New(ctx, &clientWrapper{client})
}

type clientWrapper struct {
	Client
}

func (*clientWrapper) Start(context.Context) error {
	return nil
}

func (*clientWrapper) Stop(context.Context) error {
	return nil
}
