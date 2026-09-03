package server

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metrics   *metricsCollector
	metricsMu sync.Once
)

type metricsCollector struct {
	reqCount    *prometheus.CounterVec
	reqDuration *prometheus.HistogramVec
	cacheOps    *prometheus.CounterVec
}

func getMetrics() *metricsCollector {
	metricsMu.Do(func() {
		metrics = &metricsCollector{
			reqCount: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "image_text_http_requests_total",
					Help: "Total number of HTTP requests",
				},
				[]string{"method", "path", "status"},
			),
			reqDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "image_text_http_request_duration_seconds",
					Help:    "HTTP request duration in seconds",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"method", "path"},
			),
			cacheOps: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "image_text_cache_ops_total",
					Help: "Total number of image cache operations",
				},
				[]string{"type"},
			),
		}
		prometheus.MustRegister(metrics.reqCount, metrics.reqDuration, metrics.cacheOps)
	})
	return metrics
}

func (m *metricsCollector) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		begin := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		routePath := routePattern(r)
		status := strconv.Itoa(ww.Status())
		m.reqCount.WithLabelValues(r.Method, routePath, status).Inc()
		m.reqDuration.WithLabelValues(r.Method, routePath).Observe(time.Since(begin).Seconds())
	})
}

// routePattern returns the matched chi route pattern (e.g.
// "/{size}/{bg_color}/{fg_color}") rather than the raw request path. Using
// the raw path as a metric label would create an unbounded number of time
// series, since size/color/font segments are arbitrary user input. Requests
// that don't match any route (404s, e.g. from bots probing random paths)
// are bucketed under a single "not_found" label instead.
func routePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return "unknown"
	}
	if p := rctx.RoutePattern(); p != "" {
		return p
	}
	return "not_found"
}

func (m *metricsCollector) cacheHit() {
	m.cacheOps.WithLabelValues("hit").Inc()
}

func (m *metricsCollector) cacheMiss() {
	m.cacheOps.WithLabelValues("miss").Inc()
}

func prometheusHandler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}
