package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/hanshuaikang/gin-prometheus"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"net/http"
)

const serviceName = "gin-prometheus-demo"

func main() {
	fmt.Println("initializing")

	r := gin.New()
	globalAttributes := []attribute.KeyValue{
		semconv.K8SPodName("pod-1"),
		semconv.K8SNamespaceName("test"),
		semconv.ServiceName(serviceName),
	}
	r.Use(ginprometheus.Middleware(
		// Custom attributes
		ginprometheus.WithAttributes(func(route string, request *http.Request) []attribute.KeyValue {
			attrs := []attribute.KeyValue{
				semconv.HTTPMethodKey.String(request.Method),
			}
			if route != "" {
				attrs = append(attrs, semconv.HTTPRouteKey.String(route))
			}
			return attrs
		}),
		ginprometheus.WithGlobalAttributes(globalAttributes),
		ginprometheus.WithService(serviceName, "v0.0.1"),
		ginprometheus.WithMetricPrefix("infra"),
		ginprometheus.WithPrometheusPort(4433),
		ginprometheus.WithSystemMetricDisabled(),
	))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
