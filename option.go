package ginprometheus

import (
	"go.opentelemetry.io/otel/attribute"
	"net/http"
)

// Option applies a configuration to the given config
type Option interface {
	apply(cfg *config)
}

type optionFunc func(cfg *config)

func (fn optionFunc) apply(cfg *config) {
	fn(cfg)
}

func WithAttributes(attributes func(route string, request *http.Request) []attribute.KeyValue) Option {
	return optionFunc(func(cfg *config) {
		cfg.attributes = attributes
	})
}

func WithGlobalAttributes(attributes []attribute.KeyValue) Option {
	return optionFunc(func(cfg *config) {
		globalAttributes := append(attributes[0:0], attributes...)
		cfg.globalAttributes = globalAttributes
	})
}

func WithRecordInFlightDisabled() Option {
	return optionFunc(func(cfg *config) {
		cfg.recordInFlight = false
	})
}

func WithRecordDurationDisabled() Option {
	return optionFunc(func(cfg *config) {
		cfg.recordDuration = false
	})
}

func WithRecordSizeDisabled() Option {
	return optionFunc(func(cfg *config) {
		cfg.recordSize = false
	})
}

func WithGroupedStatusDisabled() Option {
	return optionFunc(func(cfg *config) {
		cfg.groupedStatus = false
	})
}

func WithSystemMetricDisabled() Option {
	return optionFunc(func(cfg *config) {
		cfg.recordSystemMetric = false
	})
}

func WithRecorder(recorder Recorder) Option {
	return optionFunc(func(cfg *config) {
		cfg.recorder = recorder
	})
}

func WithShouldRecordFunc(shouldRecord func(route string, request *http.Request) bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.recordFilter = shouldRecord
	})
}

func WithService(serviceName string, version string) Option {
	return optionFunc(func(cfg *config) {
		cfg.serviceName = serviceName
		cfg.version = version
	})
}

func WithPrometheusPort(port int) Option {
	return optionFunc(func(cfg *config) {
		cfg.prometheusPort = port
	})
}

func WithMetricPrefix(prefix string) Option {
	return optionFunc(func(cfg *config) {
		cfg.metricPrefix = prefix
	})
}
