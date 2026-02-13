package tracer1

import (
	"context"
	"encoding/base64"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

// JaegerOption set fields
type JaegerOption func(*jaegerOptions)

type jaegerOptions struct {
	username string
	password string
}

func (o *jaegerOptions) apply(opts ...JaegerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default setting
func defaultJaegerOptions() *jaegerOptions {
	return &jaegerOptions{}
}

// WithUsername set username
func WithUsername(username string) JaegerOption {
	return func(o *jaegerOptions) {
		o.username = username
	}
}

// WithPassword set password
func WithPassword(password string) JaegerOption {
	return func(o *jaegerOptions) {
		o.password = password
	}
}

// NewJaegerExporter use jaeger collector as exporter, e.g. default url=http://localhost:14268/api/traces
func NewJaegerExporter(ctx context.Context, url string, opts ...JaegerOption) (sdkTrace.SpanExporter, error) {
	ceps := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(url),
	}

	o := defaultJaegerOptions()
	o.apply(opts...)

	if o.username != "" && o.password != "" {
		// 计算 Basic Auth 字符串
		auth := base64.StdEncoding.EncodeToString([]byte("username:password"))
		ceps = append(ceps, otlptracegrpc.WithHeaders(map[string]string{
			"Authorization": "Basic " + auth,
		}))
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(), // 如果没配 TLS
		otlptracegrpc.WithEndpoint(url),
	)
	return exporter, err
}

// NewJaegerAgentExporter use jaeger agent as exporter, e.g. host=localhost port=6831
func NewJaegerAgentExporter(ctx context.Context, host string, port string) (sdkTrace.SpanExporter, error) {
	// 默认连接到 localhost:4317 (Jaeger 的 OTLP 默认端口)
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(), // 如果没配 TLS
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", host, port)),
	)
	return exporter, err
}
