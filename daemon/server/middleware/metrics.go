package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	gometrics "github.com/docker/go-metrics"
)

var (
	metricsNS       *gometrics.Namespace
	requestDuration gometrics.LabeledTimer
	metricsOnce     sync.Once
)

type MetricsMiddleware struct {
	requestDuration gometrics.LabeledTimer
}

func NewMetricsMiddleware() *MetricsMiddleware {
	metricsOnce.Do(func() {
		metricsNS = gometrics.NewNamespace("engine", "api", nil)
		requestDuration = metricsNS.NewLabeledTimer("request_duration", "HTTP request latency", "http_method", "api_route", "api_version")
		gometrics.Register(metricsNS)
	})

	return &MetricsMiddleware{
		requestDuration: requestDuration,
	}
}

func (m *MetricsMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		// Do not track if route has not been filled.
		route := vars["route"]
		if route == "" {
			return handler(ctx, w, r, vars)
		}

		// This version is only filled in when an explicit version is requested in the path. This is intentional.
		version := vars["version"]

		start := time.Now()
		err := handler(ctx, w, r, vars)
		m.requestDuration.WithValues(r.Method, route, version).UpdateSince(start)
		return err
	}
}
