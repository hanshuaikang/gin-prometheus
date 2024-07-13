package ginprometheus

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

// httpMetricsRecorder http metric 封装了基础的http 请求指标
type httpMetricsRecorder struct {
	requestsCounter       metric.Int64UpDownCounter
	totalDuration         metric.Int64Histogram
	activeRequestsCounter metric.Int64UpDownCounter
	requestSize           metric.Int64Histogram
	responseSize          metric.Int64Histogram
	cpuUsage              metric.Float64UpDownCounter
	memoryUsage           metric.Float64UpDownCounter
}

func NewHttpMetricsRecorder(serviceName, version, metricsPrefix string) Recorder {
	metricName := func(metricName string) string {
		if len(metricsPrefix) > 0 {
			return metricsPrefix + "." + metricName
		}
		return metricName
	}
	meter := otel.Meter(serviceName, metric.WithInstrumentationVersion(version))
	requestsCounter, _ := meter.Int64UpDownCounter(metricName("http.server.request_total"), metric.WithDescription("Number of Requests"), metric.WithUnit("Count"))
	totalDuration, _ := meter.Int64Histogram(metricName("http.server.duration"), metric.WithDescription("Time Taken by request"), metric.WithUnit("Milliseconds"))
	activeRequestsCounter, _ := meter.Int64UpDownCounter(metricName("http.server.active_requests"), metric.WithDescription("Number of requests inflight"), metric.WithUnit("Count"))
	requestSize, _ := meter.Int64Histogram(metricName("http.server.request_content_length"), metric.WithDescription("Request Size"), metric.WithUnit("Bytes"))
	responseSize, _ := meter.Int64Histogram(metricName("http.server.response_content_length"), metric.WithDescription("Response Size"), metric.WithUnit("Bytes"))
	cpuUsage, _ := meter.Float64UpDownCounter(metricName("system.cpu.usage"), metric.WithDescription("CPU Usage"), metric.WithUnit("Percent"))
	memoryUsage, _ := meter.Float64UpDownCounter(metricName("system.memory.usage"), metric.WithDescription("Memory Usage"), metric.WithUnit("Percent"))

	return &httpMetricsRecorder{
		requestsCounter:       requestsCounter,
		totalDuration:         totalDuration,
		activeRequestsCounter: activeRequestsCounter,
		requestSize:           requestSize,
		responseSize:          responseSize,
		cpuUsage:              cpuUsage,
		memoryUsage:           memoryUsage,
	}
}

// AddRequests increments the number of requests being processed.
func (r *httpMetricsRecorder) AddRequests(ctx context.Context, quantity int64, attributes []attribute.KeyValue) {
	r.requestsCounter.Add(ctx, quantity, metric.WithAttributes(attributes...))
}

// ObserveHTTPRequestDuration measures the duration of an HTTP request.
func (r *httpMetricsRecorder) ObserveHTTPRequestDuration(ctx context.Context, duration time.Duration, attributes []attribute.KeyValue) {
	r.totalDuration.Record(ctx, int64(duration/time.Millisecond), metric.WithAttributes(attributes...))
}

// ObserveHTTPRequestSize measures the size of an HTTP request in bytes.
func (r *httpMetricsRecorder) ObserveHTTPRequestSize(ctx context.Context, sizeBytes int64, attributes []attribute.KeyValue) {
	r.requestSize.Record(ctx, sizeBytes, metric.WithAttributes(attributes...))
}

// ObserveHTTPResponseSize measures the size of an HTTP response in bytes.
func (r *httpMetricsRecorder) ObserveHTTPResponseSize(ctx context.Context, sizeBytes int64, attributes []attribute.KeyValue) {
	r.responseSize.Record(ctx, sizeBytes, metric.WithAttributes(attributes...))
}

// AddInflightRequests increments and decrements the number of inflight request being processed.
func (r *httpMetricsRecorder) AddInflightRequests(ctx context.Context, quantity int64, attributes []attribute.KeyValue) {
	r.activeRequestsCounter.Add(ctx, quantity, metric.WithAttributes(attributes...))
}

func (r *httpMetricsRecorder) ObserveSystemMetric(ctx context.Context, attributes []attribute.KeyValue) {
	// 这里仅作为示例，实际可能需要用更准确的方式计算 CPU 使用率
	cpuUsage := getCpuUsage()
	r.cpuUsage.Add(ctx, cpuUsage, metric.WithAttributes(attributes...))
	// 这里仅作为示例，实际可能需要用更准确的方式计算内存使用率
	memUsage := getMemoryUsage()
	r.memoryUsage.Add(ctx, memUsage, metric.WithAttributes(attributes...))
}
