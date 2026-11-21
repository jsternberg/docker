package telemetryv0

import (
	"context"

	"github.com/moby/extensions"
	"github.com/moby/extensions/clientpoint"
	tracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
)

type Client interface {
	UploadTraces(ctx context.Context, protoSpans []*tracev1.ResourceSpans) error
}

var Point = extensions.DefinePoint[Client]("org.mobyproject.extension.telemetry.collector.v0")

var ClientPoint = clientpoint.Registration{
	Point: Point.ID(),
	Provider: func(conn grpc.ClientConnInterface) extensions.Provider {
		provider := &clientProvider{
			tracesClient: tracepb.NewTraceServiceClient(conn),
		}
		return Point.Provide(provider)
	},
}

type clientProvider struct {
	tracesClient tracepb.TraceServiceClient
}

func (s *clientProvider) UploadTraces(ctx context.Context, protoSpans []*tracev1.ResourceSpans) error {
	in := &tracepb.ExportTraceServiceRequest{ResourceSpans: protoSpans}
	_, err := s.tracesClient.Export(ctx, in)
	return err
}
