package trace


import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// serviceName 服务名称
func InitTracer(endpoint, serviceName string) error {
	exp, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(200*time.Microsecond)),
		// 在资源中记录有关此应用程序的信息
		sdktrace.WithResource(
			resource.NewSchemaless(
				semconv.ServiceNameKey.String(serviceName),
				attribute.String("exporter", "otlp"),
				attribute.Float64("float", 312.23),
			),
		),
	)
	otel.SetTracerProvider(tp)

	return nil
}
