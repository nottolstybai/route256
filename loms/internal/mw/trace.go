package mw

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Trace(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	var traceID string
	if traceIDs, ok := md["x-trace-id"]; ok {
		traceID = traceIDs[0]
	}

	var spanID string
	if spanIDs, ok := md["x-span-id"]; ok {
		spanID = spanIDs[0]
	}

	if traceID != "" {
		var err error
		cfg := trace.SpanContextConfig{
			TraceFlags: trace.FlagsSampled,
			Remote:     true,
		}

		cfg.TraceID, err = trace.TraceIDFromHex(traceID)
		if err != nil {
			return nil, err
		}

		cfg.SpanID, err = trace.SpanIDFromHex(spanID)
		if err != nil {
			return nil, err
		}

		ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(cfg))
	}

	ctx, span := otel.Tracer("default").Start(ctx, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	return handler(ctx, req)
}
