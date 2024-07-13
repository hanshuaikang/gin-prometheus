package ginprometheus

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"net/http"
)

// 默认配置
type config struct {
	serviceName        string // 服务名
	version            string // 服务版本
	prometheusPort     int    //prometheus port
	metricPrefix       string // metric prefix
	recordInFlight     bool   // 是否记录当前正在处理的请求数量
	recordSize         bool   // 是否记录请求的大小
	recordDuration     bool   // 是否记录请求的处理时间
	groupedStatus      bool   // 是否对请求状态进行分
	recordSystemMetric bool   // 记录 system 指标
	recorder           Recorder
	globalAttributes   []attribute.KeyValue
	attributes         func(route string, request *http.Request) []attribute.KeyValue
	recordFilter       func(route string, request *http.Request) bool
}

func defaultConfig() *config {
	return &config{
		recordInFlight:     true,
		recordDuration:     true,
		recordSize:         true,
		groupedStatus:      true,
		recordSystemMetric: true,
		serviceName:        "gin-prometheus-service",
		version:            "v1.0.0",
		prometheusPort:     2233,
		metricPrefix:       "",
		attributes:         DefaultAttributes,
		recordFilter: func(_ string, _ *http.Request) bool {
			return true
		},
	}
}

var DefaultAttributes = func(route string, request *http.Request) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String(request.Method),
	}
	if route != "" {
		attrs = append(attrs, semconv.HTTPRouteKey.String(route))
	}
	return attrs
}
