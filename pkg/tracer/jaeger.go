package tracer

import (
	"context"
	jgr "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	tr "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"sync"
)

var (
	tracer     *JaegerTracer
	tracerOnce sync.Once
	tracerErr  error
)

type JaegerTracer struct {
	tracer tr.TracerProvider
}

func NewTracer(cfg *Config) (*JaegerTracer, error) {
	tracerOnce.Do(func() {
		if !cfg.Disabled {
			exp, err := jgr.New(jgr.WithCollectorEndpoint(jgr.WithEndpoint(cfg.URL)))
			if err != nil {
				tracerErr = err

				return
			}

			tracer = &JaegerTracer{
				tracer: trace.NewTracerProvider(
					// Always be sure to batch in production.
					trace.WithBatcher(exp),
					// Record information about this application in a Resource.
					trace.WithResource(resource.NewWithAttributes(
						semconv.SchemaURL,
						semconv.ServiceNameKey.String(cfg.ServiceName),
					))),
			}
		} else {
			tracer = &JaegerTracer{
				tracer: tr.NewNoopTracerProvider(),
			}
		}
	})

	return tracer, tracerErr
}

func (t *JaegerTracer) TracerProvider() tr.TracerProvider {
	return t.tracer
}

func (t *JaegerTracer) Start(ctx context.Context,
	tracerName, spanName string, options int64) (context.Context, tr.Span) {
	ctx, span := t.tracer.Tracer(tracerName).Start(ctx, spanName)

	if options&CtxWithTraceValue != 0 {
		ctx = context.WithValue(ctx, CtxTraceIDKey, span.SpanContext().TraceID().String())
	}

	if options&CtxWithGRPCMetadata != 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, CtxTraceIDStr, span.SpanContext().TraceID().String())
	}

	return ctx, span
}

func (t *JaegerTracer) Continue(ctx context.Context,
	tracerName, spanName string, options int64, traceID tr.TraceID) (context.Context, tr.Span) {
	spanContext := tr.NewSpanContext(tr.SpanContextConfig{
		TraceID: traceID,
	})
	// Embedding span config into the context
	ctx = tr.ContextWithSpanContext(ctx, spanContext)

	return t.Start(ctx, tracerName, spanName, options)
}

func GetTraceIDFromValue(ctx context.Context) (tr.TraceID, bool) {
	traceAny := ctx.Value(CtxTraceIDKey)

	traceStr, ok := traceAny.(string)
	if !ok {
		return tr.TraceID{}, false
	}

	return parseTraceID(traceStr)
}

func GetTraceIDFromGRPCMetadata(ctx context.Context) (tr.TraceID, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return tr.TraceID{}, false
	}

	if len(md[CtxTraceIDStr]) == 0 {
		return tr.TraceID{}, false
	}

	traceStr := md[CtxTraceIDStr][0]

	return parseTraceID(traceStr)
}

func parseTraceID(str string) (tr.TraceID, bool) {
	traceID, err := tr.TraceIDFromHex(str)
	if err != nil {
		zap.S().Error(err)

		return traceID, false
	}

	return traceID, true
}
