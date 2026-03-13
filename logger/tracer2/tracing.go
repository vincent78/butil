package tracing

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer
var JaegerServername = "test-service"

var DateTimeMilli = "2006-01-02 15:04:05.000000"

// var JaegerHost = "http://jaeger-collector.common:14268/api/traces"
var JaegerHost = "127.0.0.1:4317"

var ServerRunEnv *string // 运行环境: local、dev、test、prod

func init() {
	ServerRunEnv = new("dev")
}

// InitJaeger 初始化 Jaeger 追踪
func InitJaeger() (func(), error) {
	// 根据环境选择 Jaeger 端点
	jaegerEndpoint := JaegerHost

	// 创建 Jaeger 导出器

	// 默认连接到 localhost:4317 (Jaeger 的 OTLP 默认端口)
	exp, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithInsecure(), // 如果没配 TLS
		otlptracegrpc.WithEndpoint(jaegerEndpoint),
	)
	//exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	// 创建资源
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(JaegerServername),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String(*ServerRunEnv),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建 TracerProvider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(1.0)), // 100% 采样
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 设置全局 TextMapPropagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 创建 Tracer
	Tracer = tp.Tracer(JaegerServername)

	// 返回清理函数
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}

// StartSpan 开始一个新的 span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer.Start(ctx, name, opts...)
}

// AddSpanAttributes 添加 span 属性
func AddSpanAttributes(span trace.Span, attrs map[string]string) {
	if span == nil {
		return
	}

	var attributes []attribute.KeyValue
	for k, v := range attrs {
		attributes = append(attributes, attribute.String(k, v))
	}
	span.SetAttributes(attributes...)
}

// AddSpanEvent 添加 span 事件
func AddSpanEvent(span trace.Span, name string, attrs map[string]string) {
	if span == nil {
		return
	}

	var attributes []attribute.KeyValue
	for k, v := range attrs {
		attributes = append(attributes, attribute.String(k, v))
	}
	span.AddEvent(name, trace.WithAttributes(attributes...))
}

// StartSpanWithCustomTraceID 使用自定义（已是真实）traceID 创建新的 span
func StartSpanWithCustomTraceID(name, customTraceID string) (context.Context, trace.Span) {
	// 使用 context.Background() 作为基础 context
	ctx := context.Background()

	// customTraceID 已经是中间件生成的真实 TraceID（32位十六进制），直接解析
	var traceID trace.TraceID
	if parsed, err := trace.TraceIDFromHex(customTraceID); err == nil {
		traceID = parsed
	} else {
		// 兜底：若传入不是有效的十六进制 TraceID，则回退为哈希到固定 TraceID
		sum := sha256.Sum256([]byte(customTraceID))
		copy(traceID[:], sum[:16])
	}

	// 创建 SpanContext（仅指定 TraceID，让 OTel 生成新的 SpanID）
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     trace.SpanID{},
		TraceFlags: trace.FlagsSampled,
	})

	// 将 SpanContext 注入到 context 中
	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	// 创建 span
	ctx, span := Tracer.Start(ctx, name)

	// 添加自定义 traceID 作为属性
	span.SetAttributes(attribute.String("custom_trace_id", customTraceID))

	return ctx, span
}

func StartTimedSpan(ctx context.Context, name string, attrs map[string]string) func() {
	_, span := StartSpan(ctx, name)
	start := time.Now()

	return func() {
		end := time.Now()

		spanAttrs := make(map[string]string, len(attrs)+3)
		for k, v := range attrs {
			spanAttrs[k] = v
		}
		spanAttrs["send_start_time"] = start.Format(DateTimeMilli)
		spanAttrs["send_end_time"] = end.Format(DateTimeMilli)
		spanAttrs["send_duration_ms"] = strconv.FormatInt(end.Sub(start).Milliseconds(), 10)

		AddSpanAttributes(span, spanAttrs)
		span.End()
	}
}
