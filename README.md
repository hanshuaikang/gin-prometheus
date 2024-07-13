<h1 align="center">🛰 gin-prometheus</h1>
<p align="center">
    <em>Prometheus metrics exporter for Gin base on OpenTelemetry.</em>
</p>

### 🔰 Installation

```shell
$ go get -u github.com/hanshuaikang/gin-prometheus
```

### 📝 Usage

It's easy to get started with in-prometheus, You just need to install a middleware for your Gin project

```golang
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/hanshuaikang/gin-prometheus"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"net/http"
)

const serviceName = "gin-prometheus-demo"

func main() {
	fmt.Println("initializing")
	r := gin.New()
	r.Use(ginprometheus.Middleware())
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
```

The Prometheus service will start, you can visit http://localhost:2233/metrics at your browser

gin-prometheus supports a lot of custom configuration. the detail please look example/example.go file


### 🎉 Metrics

Details about exposed Prometheus metrics.

| Name                                   | Type | Exposed Information           |
|----------------------------------------| ---- |-------------------------------|
| http_server_active_requests						      | Counter	| Number of requests inflight.  |
| http_server_duration		                 | Histogram	| Time Taken by request         |
| http_server_request_total              | Counter | Number of Requests.           |
| http_server_request_content_length 		  | Histogram	| HTTP request sizes in bytes.  |
| http_server_response_content_length 		 | Histogram	| HTTP response sizes in bytes. |
| system_cpu_usage 		                    | Counter	| CPU Usage                     |
| system_memory_usage		                  | Counter	| Memory Usage                  |
