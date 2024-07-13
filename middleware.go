package ginprometheus

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"time"
)

func recordSystemMetrics(ctx context.Context, recorder Recorder, interval time.Duration, attributes []attribute.KeyValue) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			recorder.ObserveSystemMetric(ctx, attributes)
		case <-ctx.Done():
			// 上下文取消信号，退出循环
			return
		}
	}
}

// Middleware 请求中间件
func Middleware(options ...Option) gin.HandlerFunc {

	ctx := context.Background()
	// 默认配置
	cfg := defaultConfig()
	for _, option := range options {
		option.apply(cfg)
	}

	// start the prometheus server
	initMetrics(cfg.prometheusPort, cfg.serviceName)

	// 初始化记录器
	recorder := cfg.recorder
	if recorder == nil {
		recorder = NewHttpMetricsRecorder(cfg.serviceName, cfg.version, cfg.metricPrefix)
	}

	if cfg.recordSystemMetric {
		go recordSystemMetrics(ctx, recorder, time.Minute, cfg.globalAttributes)
	}

	return func(ginCtx *gin.Context) {
		// 获取这次请求的完整路径
		route := ginCtx.FullPath()
		if len(route) <= 0 {
			ginCtx.Next()
			return
		}

		if !cfg.recordFilter(route, ginCtx.Request) {
			ginCtx.Next()
			return
		}

		start := time.Now()
		reqAttributes := append(cfg.attributes(route, ginCtx.Request), cfg.globalAttributes...)

		if cfg.recordInFlight {
			// count 类型，需要开始时 + 1, 结束时 -1
			recorder.AddInflightRequests(ctx, 1, reqAttributes)
			defer recorder.AddInflightRequests(ctx, -1, reqAttributes)
		}

		defer func() {
			// generate a new slice
			resAttributes := append(reqAttributes[0:0], reqAttributes...)
			if cfg.groupedStatus {
				// 200 300 400 500
				code := int(ginCtx.Writer.Status()/100) * 100
				resAttributes = append(resAttributes, semconv.HTTPStatusCodeKey.Int(code))
			} else {
				resAttributes = append(resAttributes, semconv.HTTPStatusCodeKey.Int(ginCtx.Writer.Status()))
			}

			recorder.AddRequests(ctx, 1, resAttributes)

			if cfg.recordSize {
				requestSize := computeApproximateRequestSize(ginCtx.Request)
				recorder.ObserveHTTPRequestSize(ctx, requestSize, resAttributes)
				recorder.ObserveHTTPResponseSize(ctx, int64(ginCtx.Writer.Size()), resAttributes)
			}

			if cfg.recordDuration {
				recorder.ObserveHTTPRequestDuration(ctx, time.Since(start), resAttributes)
			}
		}()

		ginCtx.Next()
	}
}
